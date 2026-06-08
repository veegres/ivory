package node

import (
	"ivory/clients/console"
	"ivory/core/service/job"
	"ivory/plugins/platform"
)

func (s *Service) PlatformCopyId(r PlatformCopyIdRequest) (string, error) {
	adapter, err := s.platformRegistry.Get(platform.Onprem)
	if err != nil {
		return "", err
	}
	con := s.getPlatformCredConnection(r.PlatformCredConnection)
	return "ok", adapter.CopyId(con, r.PublicKey)
}

func (s *Service) PlatformMetrics(r PlatformMetricsRequest) (*PlatformMetricsResponse, error) {
	adapter, conn, err := s.getPlatformAdapter(r)
	if err != nil {
		return nil, err
	}
	return adapter.Metrics(conn)
}

func (s *Service) PlatformContainerStop(r PlatformActionRequest) ([]string, error) {
	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(adapter.ListContainer(conn))
}

func (s *Service) PlatformContainerStart(r PlatformActionRequest) ([]string, error) {
	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(adapter.ListContainer(conn))
}

func (s *Service) PlatformContainerUp(r PlatformUpRequest, subscriberID job.SubscriberID, send func(event job.Message)) {
	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		send(job.Message{Type: job.SERVER, Message: err.Error()})
		return
	}
	s.streamCommand(adapter.UpContainer(conn, r.Options, r.Image), subscriberID, send)
}

func (s *Service) PlatformContainerDown(r PlatformActionRequest) ([]string, error) {
	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(adapter.ListContainer(conn))
}

func (s *Service) PlatformContainerList(c PlatformVaultConnection) ([]string, error) {
	adapter, conn, err := s.getPlatformAdapter(c)
	if err != nil {
		return nil, err
	}
	return s.executeCommand(adapter.ListContainer(conn))
}

func (s *Service) PlatformContainerLogs(r PlatformLogsRequest, subscriberID job.SubscriberID, send func(event job.Message)) {
	adapter, conn, err := s.getPlatformAdapter(r.Connection)
	if err != nil {
		send(job.Message{Type: job.SERVER, Message: err.Error()})
		return
	}
	s.streamCommand(adapter.LogsContainer(conn, r.Name, r.Tail), subscriberID, send)
}

func (s *Service) streamCommand(cmd console.Command, subscriberID job.SubscriberID, send func(event job.Message)) {
	jobID, err := s.jobManager.Start(cmd)
	if err != nil {
		send(job.Message{Type: job.SERVER, Message: err.Error()})
		return
	}
	s.jobManager.Stream(jobID, subscriberID, send)
}

func (s *Service) executeCommand(cmd console.Command) ([]string, error) {
	return cmd.Execute()
}
