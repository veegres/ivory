// Package platform defines the plugin boundary for deployment targets
// (on-prem Docker-over-SSH now, Kubernetes/OpenShift later). A
// platform.Adapter lets the rest of Ivory read host/container metrics and
// deploy/manage containers without knowing which concrete platform it is
// talking to.
package platform

import (
	"errors"
	"ivory/clients/console"
)

var ErrInvalidCpuMetrics = errors.New("invalid cpu metrics output")
var ErrInvalidMemoryMetrics = errors.New("invalid memory metrics output")
var ErrInvalidNetworkMetrics = errors.New("invalid network metrics output")
var ErrInvalidContainerMetrics = errors.New("invalid container metrics output")
var ErrInvalidProcesses = errors.New("invalid processes output")

// Connection contains the bare minimum details to execute commands on a
// platform node. PrivateKey takes precedence over Password when both are set.
type Connection struct {
	Host       string
	Port       int
	Username   string
	Password   string
	PrivateKey []byte
}

// Adapter is implemented by every platform plugin (linux, ...). It combines
// host-level (VmManager) and container-level (ContainerManager) operations
// into the single interface the rest of Ivory depends on.
type Adapter interface {
	VmManager
	ContainerManager
}

// VmManager covers host-level operations against the node itself, not any
// specific container running on it.
type VmManager interface {
	// Metrics returns the node's current CPU/memory/network usage.
	Metrics(connection Connection) (*Metrics, error)
	// CopyId installs publicKey into the node's authorized_keys, enabling
	// key-based SSH access.
	CopyId(connection Connection, publicKey string) error
	// Logs streams (or dumps, if follow is false) the file at path.
	Logs(connection Connection, path string, tail int, follow bool) console.Command
	// Processes lists the node's running processes.
	Processes(connection Connection) ([]Process, error)
	// Info returns a list of arbitrary platform-reported details about the
	// node (OS, kernel, uptime, hardware...). Content and order are entirely
	// up to each adapter.
	Info(connection Connection) ([]InfoItem, error)
}

// ContainerManager covers the lifecycle and inspection of a single deployed
// container on the node.
type ContainerManager interface {
	// RenderOptions renders a deploy spec into the adapter's native, still
	// user-editable options text (docker run flags, a manifest, etc).
	RenderOptions(spec DeploySpec) string
	// ListContainer lists every deployed container on the node.
	ListContainer(connection Connection) console.Command
	// UpContainer creates and starts a new container running image with the
	// given native options text (as rendered by RenderOptions). entryScript,
	// if non-empty, replaces the image's own default startup command (see
	// keeper.DeploymentSpec.EntryScript).
	UpContainer(connection Connection, options, image, entryScript string) console.Command
	// DownContainer removes the named container.
	DownContainer(connection Connection, name string) console.Command
	// StartContainer starts an existing, stopped container.
	StartContainer(connection Connection, name string) console.Command
	// StopContainer stops a running container without removing it.
	StopContainer(connection Connection, name string) console.Command
	// RestartContainer restarts a container.
	RestartContainer(connection Connection, name string) console.Command
	// ExecContainer runs a command inside the named running container.
	ExecContainer(connection Connection, name string, command string) console.Command
	// LogsContainer streams (or dumps, if follow is false) the named
	// container's logs.
	LogsContainer(connection Connection, name string, tail int, follow bool) console.Command
	// MetricsContainer returns the named container's current CPU/memory/network usage.
	MetricsContainer(connection Connection, name string) (*Metrics, error)
}
