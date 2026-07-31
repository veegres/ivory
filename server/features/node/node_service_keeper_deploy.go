package node

import (
	"errors"
	"fmt"
	"ivory/plugins/keeper"
	"ivory/plugins/platform"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var ErrKeeperDeployImageNotProvided = errors.New("deploy image not provided")
var ErrKeeperDeployDatabaseCredentialsRequired = errors.New("database credentials are required")
var ErrKeeperDeployPlanHasWarnings = errors.New("deployment plan has warnings")

func (s *Service) KeeperDeploySpec(r KeeperDeploySpecRequest) (*KeeperDeploySpecResponse, error) {
	metadata, err := s.keeperMetadataRegistry.Get(r.Plugin)
	if err != nil {
		return nil, err
	}
	spec := metadata.DeploymentSpec()
	return &KeeperDeploySpecResponse{
		Uri:    spec.DefaultImage,
		Fields: mapKeeperDeploymentFields(spec),
	}, nil
}

// KeeperDeployPlan resolves a deployment intent into concrete
// per-node deployments: it fills port defaults from the keeper plugin spec,
// resolves every declared field (a user-provided value wins, otherwise the
// value is derived from the node list or falls back to the default; port
// fields are unique per node in single-host mode), renders the options
// template unless a node carries an override, and interpolates a
// credential-masked preview. It is a pure computation: no vaults are read,
// nothing is executed.
func (s *Service) KeeperDeployPlan(r KeeperDeployPlanRequest) (*KeeperDeployPlanResponse, error) {
	metadata, err := s.keeperMetadataRegistry.Get(r.Plugin)
	if err != nil {
		return nil, err
	}
	adapter, err := s.platformRegistry.Get(platform.Linux)
	if err != nil {
		return nil, err
	}
	spec := metadata.DeploymentSpec()

	image := r.Image
	if image == "" {
		image = spec.DefaultImage
	}

	warnings := make([]string, 0)
	template := adapter.RenderOptions(mapKeeperDeploymentToPlatformSpec(spec, r.SingleHost))

	dbPortDefault, _ := strconv.Atoi(spec.Defaults[keeper.VarDbPort])
	keeperPortDefault, hasKeeperEndpoint := spec.Defaults[keeper.VarKeeperPort]

	// NOTE: EntryScript applies to every node by default, the first node
	// (the leader; see keeper.DeploymentSpec.EntryScript's doc for why it's
	// first) included, unless the plugin sets EntryScriptReplicasOnly, in
	// which case the leader gets the image's normal command instead and
	// only the other nodes get EntryScript. Either way, VarLeaderHost is
	// resolved to the leader's real host, already known from the request
	// itself.
	var leaderHost string
	if len(r.Nodes) > 0 {
		leaderHost = r.Nodes[0].Host
	}

	nodes := make([]KeeperDeployPlanNode, 0, len(r.Nodes))
	for i, n := range r.Nodes {
		dbPort := dbPortDefault
		if n.DbPort != nil && *n.DbPort > 0 {
			dbPort = *n.DbPort
		} else if r.SingleHost {
			// NOTE: mirrors the FieldPort base+i offset below - nodes sharing
			// one host can't all bind the same default database port
			dbPort += i
		}
		keeperPort := dbPort
		if hasKeeperEndpoint {
			keeperPort, _ = strconv.Atoi(keeperPortDefault)
			if n.KeeperPort != nil && *n.KeeperPort > 0 {
				keeperPort = *n.KeeperPort
			} else if r.SingleHost {
				keeperPort += i
			}
		} else if n.KeeperPort != nil && *n.KeeperPort > 0 && *n.KeeperPort != dbPort {
			warnings = append(warnings, fmt.Sprintf(
				"host %q: keeper port %d is ignored, the plugin serves its keeper endpoint on the database port %d",
				n.Host, *n.KeeperPort, dbPort,
			))
			// NOTE: honor the plugin convention keeperPort == dbPort
		}
		sshPort := 22
		if n.SshPort != nil && *n.SshPort > 0 {
			sshPort = *n.SshPort
		}
		options := n.Options
		if options == "" {
			options = template
		}
		var entryScript string
		if (i > 0 || !spec.EntryScriptReplicasOnly) && spec.EntryScript != "" {
			entryScript = strings.ReplaceAll(spec.EntryScript, string(keeper.VarLeaderHost), leaderHost)
		}
		nodes = append(nodes, KeeperDeployPlanNode{
			Host:        n.Host,
			SshPort:     sshPort,
			KeeperPort:  keeperPort,
			DbPort:      dbPort,
			Ports:       make(map[string]int, 1),
			Options:     options,
			EntryScript: entryScript,
		})
	}

	effective := make(map[string]string, len(spec.Fields))

	// 1. port fields: the effective value is the base, every node gets its
	// own value in single-host mode so the listeners don't collide
	for _, f := range spec.Fields {
		if f.Type != keeper.FieldPort {
			continue
		}
		raw := r.Values[string(f.Name)]
		if raw == "" {
			raw = f.Default
		}
		base, errPort := strconv.Atoi(raw)
		if errPort != nil || base <= 0 {
			warnings = append(warnings, fmt.Sprintf("field %q: invalid port %q, using the default %s", f.Name, raw, f.Default))
			base, _ = strconv.Atoi(f.Default)
		}
		effective[string(f.Name)] = strconv.Itoa(base)
		for i := range nodes {
			if r.SingleHost {
				nodes[i].Ports[string(f.Name)] = base + i
			} else {
				nodes[i].Ports[string(f.Name)] = base
			}
		}
	}

	// 2. text fields in declaration order, so a template may reference
	// earlier fields (e.g. the member list references the peer port)
	for _, f := range spec.Fields {
		if f.Type == keeper.FieldPort {
			continue
		}
		if value := r.Values[string(f.Name)]; value != "" {
			effective[string(f.Name)] = value
			continue
		}
		if f.Template != "" {
			entries := make([]string, 0, len(nodes))
			for i := range nodes {
				entries = append(entries, keeper.Interpolate(f.Template, s.getPlanValues(r, nodes[i], effective)))
			}
			effective[string(f.Name)] = strings.Join(entries, f.Separator)
			continue
		}
		effective[string(f.Name)] = f.Default
	}

	_, hasCredentials := spec.Defaults[keeper.VarDbUser]
	for i := range nodes {
		nodes[i].OptionsPreview = keeper.Interpolate(nodes[i].Options, s.getPlanValues(r, nodes[i], effective))
		nodes[i].EntryScriptPreview = keeper.Interpolate(nodes[i].EntryScript, s.getPlanValues(r, nodes[i], effective))
		unresolved := keeper.UnresolvedPlaceholders(nodes[i].OptionsPreview)
		unresolved = append(unresolved, keeper.UnresolvedPlaceholders(nodes[i].EntryScriptPreview)...)
		for _, u := range unresolved {
			// NOTE: credentials are resolved from the vault only at execution
			// time, so the preview intentionally keeps their variables
			// visible instead of faking values, and they are not missing
			if hasCredentials && (u == string(keeper.VarDbUser) || u == string(keeper.VarDbPass)) {
				continue
			}
			warning := fmt.Sprintf("missing value for placeholder %s", u)
			if !slices.Contains(warnings, warning) {
				warnings = append(warnings, warning)
			}
		}
	}

	return &KeeperDeployPlanResponse{
		Image:      image,
		Values:     effective,
		PostScript: spec.PostScript,
		Fields:     mapKeeperDeploymentFields(spec),
		Nodes:      nodes,
		Warnings:   warnings,
	}, nil
}

// getPlanValues assembles the interpolation values for one planned node from
// the request values, the node's resolved ports and the effective field
// values; effective values win over raw request ones and per-node port
// values win over the cluster-level base.
func (s *Service) getPlanValues(r KeeperDeployPlanRequest, n KeeperDeployPlanNode, effective map[string]string) map[string]string {
	values := make(map[string]string, len(r.Values)+len(effective)+len(n.Ports)+4)
	maps.Copy(values, r.Values)
	maps.Copy(values, effective)
	values[string(keeper.VarCluster)] = r.Cluster
	values[string(keeper.VarHost)] = n.Host
	values[string(keeper.VarKeeperPort)] = strconv.Itoa(n.KeeperPort)
	values[string(keeper.VarDbPort)] = strconv.Itoa(n.DbPort)
	for name, port := range n.Ports {
		values[name] = strconv.Itoa(port)
	}
	return values
}

// buildNodeValues assembles the interpolation values for one node from the
// request values and the plan's effective field values; credentials are
// stripped because they are resolved from the vault at execution time,
// effective field values win over raw request-supplied ones and per-node
// port values win over the base.
func (s *Service) buildNodeValues(name string, requestValues map[string]string, planValues map[string]string, pn KeeperDeployPlanNode) map[string]string {
	values := make(map[string]string, len(requestValues)+len(planValues)+len(pn.Ports)+3)
	maps.Copy(values, requestValues)
	maps.Copy(values, planValues)
	delete(values, string(keeper.VarDbUser))
	delete(values, string(keeper.VarDbPass))
	values[string(keeper.VarCluster)] = name
	values[string(keeper.VarKeeperPort)] = strconv.Itoa(pn.KeeperPort)
	values[string(keeper.VarDbPort)] = strconv.Itoa(pn.DbPort)
	for portName, port := range pn.Ports {
		values[portName] = strconv.Itoa(port)
	}
	return values
}

// KeeperDeployUp deploys one node that a KeeperDeployPlan already
// resolved (ports, options): it builds the node's interpolation values and
// executes the generic container deploy primitive.
func (s *Service) KeeperDeployUp(r KeeperDeployUpRequest) ([]string, error) {
	if r.Node.Host == "" {
		return nil, errors.New("host not provided for node")
	}
	values := s.buildNodeValues(r.Cluster, r.RequestValues, r.PlanValues, r.Node)
	return s.PlatformContainerUp(PlatformUpRequest{
		// NOTE: the deployed container's name is the node's own host, matching
		// what KeeperPostDeploy/PlatformContainerStop/Start/Exec assume when
		// they address it by r.Node.Host afterwards - it must not be the
		// cluster name, which every node in a multi-node deploy would share.
		Name:        r.Node.Host,
		Image:       r.Image,
		Connection:  r.Connection,
		Vaults:      r.Vaults,
		Options:     r.Node.Options,
		EntryScript: r.Node.EntryScript,
		Values:      values,
	})
}

// KeeperPostDeploy runs a deployment plan's post-deploy script (e.g.
// enabling authentication) once inside one already-deployed node. A failing
// script is reported as a log line rather than an error, since the
// deployment itself already succeeded.
func (s *Service) KeeperPostDeploy(r KeeperPostDeployRequest) []string {
	if r.PostScript == "" {
		return nil
	}
	logs := make([]string, 0, 1)

	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		logs = append(logs, fmt.Sprintf("post-deploy initialization failed: %v", err))
		return logs
	}
	nodeValues := s.buildNodeValues(r.Cluster, r.RequestValues, r.PlanValues, r.Node)
	values, err := s.getExecutionValues(r.Connection.Host, r.Vaults, nodeValues)
	if err != nil {
		logs = append(logs, fmt.Sprintf("post-deploy initialization failed: %v", err))
		return logs
	}

	res, err := s.execContainerCommand(adapter, conn, r.Node.Host, r.PostScript, values)
	if err != nil {
		logs = append(logs, fmt.Sprintf("post-deploy initialization failed: %v", err))
		return logs
	}
	logs = append(logs, res...)
	return logs
}

// KeeperDeploy resolves and deploys a single keeper node end-to-end: it plans
// the deployment for this one node (ports, options, interpolation), executes
// the container deploy, and runs any post-deploy initialization the keeper
// plugin declares. It is the self-contained deploy action for one node; a
// multi-node cluster deploy instead calls KeeperDeployUp per node and
// KeeperPostDeploy once, since post-deploy initialization runs only once per
// cluster, not once per node.
func (s *Service) KeeperDeploy(r KeeperDeployRequest) ([]string, error) {
	plan, err := s.KeeperDeployPlan(KeeperDeployPlanRequest{
		Plugin:     r.Plugin,
		Cluster:    r.Cluster,
		SingleHost: r.SingleHost,
		Image:      r.Image,
		Values:     r.Values,
		Nodes:      []KeeperDeployPlanNodeRequest{r.Node},
	})
	if err != nil {
		return nil, err
	}
	if plan.Image == "" {
		return nil, ErrKeeperDeployImageNotProvided
	}
	if len(plan.Warnings) > 0 {
		return nil, fmt.Errorf("%w: %s", ErrKeeperDeployPlanHasWarnings, strings.Join(plan.Warnings, "; "))
	}
	if _, dbCredentials := plan.Fields.Defaults[string(keeper.VarDbUser)]; dbCredentials && r.Vaults.DatabaseId == uuid.Nil {
		return nil, ErrKeeperDeployDatabaseCredentialsRequired
	}
	if err := s.ValidateKeeperLockedCredentials(plan.Fields.Defaults[string(keeper.VarDbUser)], r.Vaults.DatabaseId); err != nil {
		return nil, err
	}

	planNode := plan.Nodes[0]
	execRequest := KeeperDeployExecRequest{
		Cluster:       r.Cluster,
		PlanValues:    plan.Values,
		RequestValues: r.Values,
		Node:          planNode,
		Connection:    r.Connection,
		Vaults:        r.Vaults,
	}
	logs, err := s.KeeperDeployUp(KeeperDeployUpRequest{KeeperDeployExecRequest: execRequest, Image: plan.Image})
	if err != nil {
		return nil, err
	}

	if plan.PostScript != "" {
		logs = append(logs, s.KeeperPostDeploy(KeeperPostDeployRequest{KeeperDeployExecRequest: execRequest, PostScript: plan.PostScript})...)
	}
	return logs, nil
}

// ValidateKeeperLockedCredentials rejects a database vault whose username
// differs from the engine-required (locked) one declared by the keeper
// plugin's DeploymentSpec, so a deployment can never end up with credentials
// its engine would refuse. A uuid.Nil databaseId means no vault is linked yet.
func (s *Service) ValidateKeeperLockedCredentials(requiredUser string, databaseId uuid.UUID) error {
	if requiredUser == "" || databaseId == uuid.Nil {
		return nil
	}
	dbVault, err := s.vaultService.Get(databaseId)
	if err != nil {
		return err
	}
	if dbVault.Username != requiredUser {
		return fmt.Errorf("database vault username %q is not allowed: the keeper plugin locks it to %q", dbVault.Username, requiredUser)
	}
	return nil
}
