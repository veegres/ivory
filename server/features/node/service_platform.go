package node

import (
	"ivory/features/job"
	"ivory/plugins/platform"
)

func (s *Service) PlatformCopyId(c PlatformCredConnection, publicKey string) error {
	adapter, err := s.platformRegistry.Get(platform.Onprem)
	if err != nil {
		return err
	}
	con := s.getPlatformCredConnection(c)
	return adapter.CopyId(con, publicKey)
}

func (s *Service) PlatformMetrics(request PlatformVaultConnection) (*platform.Metrics, error) {
	adapter, conn, err := s.getPlatformAdapter(request)
	if err != nil {
		return nil, err
	}
	return adapter.Metrics(conn)
}

func (s *Service) PlatformStop(request PlatformDeployRequest, subscriberID job.SubscriberID, send func(event job.Message)) {
	adapter, conn, err := s.getPlatformAdapter(request.Connection)
	if err != nil {
		send(job.Message{Type: job.SERVER, Message: err.Error()})
		return
	}
	s.streamCommand(adapter.StopCommand(conn, request.Name), subscriberID, send)
}

func (s *Service) PlatformDeploy(request PlatformDeployRequest, subscriberID job.SubscriberID, send func(event job.Message)) {
	adapter, conn, err := s.getPlatformAdapter(request.Connection)
	if err != nil {
		send(job.Message{Type: job.SERVER, Message: err.Error()})
		return
	}
	s.streamCommand(adapter.DeployCommand(conn, request.Options, request.Image), subscriberID, send)
}

func (s *Service) PlatformDelete(request PlatformDeployRequest, subscriberID job.SubscriberID, send func(event job.Message)) {
	adapter, conn, err := s.getPlatformAdapter(request.Connection)
	if err != nil {
		send(job.Message{Type: job.SERVER, Message: err.Error()})
		return
	}
	s.streamCommand(adapter.DeleteCommand(conn, request.Name), subscriberID, send)
}

func (s *Service) PlatformList(c PlatformVaultConnection, subscriberID job.SubscriberID, send func(event job.Message)) {
	adapter, conn, err := s.getPlatformAdapter(c)
	if err != nil {
		send(job.Message{Type: job.SERVER, Message: err.Error()})
		return
	}
	s.streamCommand(adapter.ListCommand(conn), subscriberID, send)
}

func (s *Service) PlatformLogs(request PlatformLogsRequest, subscriberID job.SubscriberID, send func(event job.Message)) {
	adapter, conn, err := s.getPlatformAdapter(request.Connection)
	if err != nil {
		send(job.Message{Type: job.SERVER, Message: err.Error()})
		return
	}
	s.streamCommand(adapter.LogsCommand(conn, request.Name, request.Tail), subscriberID, send)
}

func (s *Service) streamCommand(cmd job.Command, subscriberID job.SubscriberID, send func(event job.Message)) {
	jobID, err := s.jobManager.Start(cmd)
	if err != nil {
		send(job.Message{Type: job.SERVER, Message: err.Error()})
		return
	}
	s.jobManager.Stream(jobID, subscriberID, send)
}
