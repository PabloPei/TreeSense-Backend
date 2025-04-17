package audit

type ActivityLog struct {
    UserID []uint8 `json:"user_id"`
    Action string  `json:"action"`
}

type AuditService interface {
	LogActivity(log ActivityLog) error
}

type AuditRepository interface {
	LogActivity(userID []uint8, action string) error
}

