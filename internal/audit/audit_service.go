package audit

type Service struct {
	repository AuditRepository
}

func NewService(repository AuditRepository) *Service {
	return &Service{ repository: repository}
}

func (s *Service) LogActivity(log ActivityLog) error {
	return s.repository.LogActivity(log.UserID, log.Action)
}