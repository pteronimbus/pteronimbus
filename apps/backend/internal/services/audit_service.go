package services

import "log"

// AuditService handles audit logging.
type AuditService struct{}

// NewAuditService creates a new AuditService.
func NewAuditService() *AuditService {
	return &AuditService{}
}

// Log logs an audit event.
func (s *AuditService) Log(event string, details map[string]interface{}) {
	log.Printf("[AUDIT] %s: %v\n", event, details)
}
