package pg_compacttable

import (
	"errors"
	"ivory/clients/console/shell"
	"ivory/core/config"
	"ivory/core/service/job"
	"ivory/core/service/vault"
	"ivory/plugins/database"

	"github.com/google/uuid"
)

var ErrJobIsActive = errors.New("job is active")
var ErrNoSuchActiveJob = errors.New("there is no such active job")

type Service struct {
	bloatRepository *Repository
	vaultService    *vault.Service
	shellClient     *shell.Client
	jobManager      *job.Service
}

func NewService(
	bloatRepository *Repository,
	shellClient *shell.Client,
	vaultService *vault.Service,
	jobManager *job.Service,
) *Service {
	service := &Service{
		bloatRepository: bloatRepository,
		vaultService:    vaultService,
		shellClient:     shellClient,
		jobManager:      jobManager,
	}

	go service.initializer()
	return service
}

func (s *Service) initializer() {
	runningJobs, _ := s.bloatRepository.ListByStatus(job.RUNNING)
	for _, runningJob := range runningJobs {
		_ = s.bloatRepository.UpdateStatus(runningJob, job.FAILED)
	}
	pendingJobs, _ := s.bloatRepository.ListByStatus(job.PENDING)
	for _, pendingJob := range pendingJobs {
		go s.start(&pendingJob)
	}
}

func (s *Service) SupportedFeatures(plugin env.Plugin) []env.Feature {
	// NOTE: pgcompacttable works only with postgres (pgstattuple based)
	if plugin.String() != database.POSTGRES.String() {
		return []env.Feature{}
	}
	return []env.Feature{
		env.ViewToolPgCompactTableList,
		env.ViewToolPgCompactTableItem,
		env.ViewToolPgCompactTableLogs,
		env.ManageToolPgCompactTableJob,
	}
}

func (s *Service) List() ([]Response, error) {
	return s.bloatRepository.List()
}

func (s *Service) ListByStatus(status job.Status) ([]Response, error) {
	return s.bloatRepository.ListByStatus(status)
}

func (s *Service) ListByCluster(cluster string) ([]Response, error) {
	return s.bloatRepository.ListByCluster(cluster)
}

func (s *Service) Get(uuid uuid.UUID) (Response, error) {
	return s.bloatRepository.Get(uuid)
}

func (s *Service) GetLogsPath(uuid uuid.UUID) (string, error) {
	model, err := s.bloatRepository.Get(uuid)
	if err != nil {
		return "", err
	}
	return s.jobManager.GetLogsPath(model.JobId)
}

func (s *Service) Start(clusterName string, vaultId *uuid.UUID, args []string) (*Response, error) {
	compactTable, err := s.bloatRepository.Create(clusterName, vaultId, args)
	if err != nil {
		return nil, err
	}
	s.start(compactTable)
	return compactTable, nil
}

func (s *Service) Stream(jobUuid uuid.UUID, subscriberID job.SubscriberID, close <-chan struct{}, send func(event job.Message)) {
	model, errModel := s.bloatRepository.Get(jobUuid)
	if errModel != nil {
		send(job.Message{Type: job.SERVER, Message: "Streaming failed: Stream Not Found"})
		send(job.Message{Type: job.STATUS, Message: job.UNKNOWN.String()})
		return
	}

	send(job.Message{Type: job.STATUS, Message: model.Status.String()})

	s.jobManager.Stream(model.JobId, subscriberID, close, send)

	latest, errGet := s.bloatRepository.Get(jobUuid)
	if errGet == nil {
		send(job.Message{Type: job.STATUS, Message: latest.Status.String()})
	} else {
		send(job.Message{Type: job.STATUS, Message: job.UNKNOWN.String()})
	}
}

func (s *Service) Delete(jobUuid uuid.UUID) error {
	model, err := s.bloatRepository.Get(jobUuid)
	if err == nil && s.isActive(model.Status) {
		return ErrJobIsActive
	}
	return s.bloatRepository.Delete(jobUuid)
}

func (s *Service) DeleteAll() error {
	runningJobs, _ := s.bloatRepository.ListByStatus(job.RUNNING)
	pendingJobs, _ := s.bloatRepository.ListByStatus(job.PENDING)
	for _, model := range append(runningJobs, pendingJobs...) {
		if model.JobId != "" {
			_ = s.jobManager.Stop(model.JobId)
		}
		_ = s.bloatRepository.UpdateStatus(model, job.FAILED)
	}
	return s.bloatRepository.DeleteAll()
}

func (s *Service) Stop(jobUuid uuid.UUID) error {
	model, err := s.bloatRepository.Get(jobUuid)
	if err != nil || !s.isActive(model.Status) || model.JobId == "" {
		return ErrNoSuchActiveJob
	}
	return s.jobManager.Stop(model.JobId)
}

func (s *Service) start(model *Response) {
	args, errArgs := s.getArgs(model)
	if errArgs != nil {
		_ = s.bloatRepository.UpdateStatus(*model, job.FAILED)
		return
	}
	command := s.shellClient.Command("pgcompacttable", args)
	command.JobID = string(model.JobId)
	command.JobPersist = true
	command.JobKeepAlive = true

	jobID, errStart := s.jobManager.Start(command)
	if errStart != nil {
		_ = s.bloatRepository.UpdateStatus(*model, job.FAILED)
		return
	}

	go s.updateJobStatus(model.Uuid, jobID)
}

func (s *Service) updateJobStatus(jobUuid uuid.UUID, jobID job.JobID) {
	subscriberID := job.SubscriberID("pg_compacttable_updater")

	channel, errSubscribe := s.jobManager.Subscribe(jobID, subscriberID)
	if errSubscribe != nil || channel == nil {
		return
	}

	for event := range channel {
		model, errGet := s.bloatRepository.Get(jobUuid)
		if errGet != nil {
			continue
		}

		switch event.Type {
		case job.STATUS:
			status := s.parseStatus(event.Message)
			_ = s.bloatRepository.UpdateStatus(model, status)
		default:
		}
	}

	s.jobManager.Unsubscribe(jobID, subscriberID)
}

func (s *Service) parseStatus(value string) job.Status {
	switch value {
	case job.PENDING.String():
		return job.PENDING
	case job.RUNNING.String():
		return job.RUNNING
	case job.FINISHED.String():
		return job.FINISHED
	case job.FAILED.String():
		return job.FAILED
	case job.STOPPED.String():
		return job.STOPPED
	default:
		return job.UNKNOWN
	}
}

func (s *Service) isActive(status job.Status) bool {
	return status == job.PENDING || status == job.RUNNING
}

func (s *Service) getArgs(model *Response) ([]string, error) {
	var vaultArgs []string
	if model.VaultId == nil {
		vaultArgs = []string{"--user", "postgres"}
	} else {
		cred, errCred := s.vaultService.GetDecrypted(*model.VaultId)
		if errCred != nil {
			return nil, errCred
		}
		vaultArgs = []string{"--user", cred.Username, "--password", cred.Secret}
	}

	args := make([]string, 0, len(model.CommandArgs)+len(vaultArgs))
	args = append(args, model.CommandArgs...)
	args = append(args, vaultArgs...)
	return args, nil
}
