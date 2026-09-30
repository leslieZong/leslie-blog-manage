package dto

import "leslie-blog-server/internal/modules/audit/model"

type CreateAuditLogRequest struct {
	UserID       *string
	Action       string
	Resource     string
	ResourceID   *string
	RequestID    string
	IP           string
	UserAgent    string
	Result       string
	ErrorMessage *string
}

type AuditLogResponse struct {
	ID string `json:"id"`

	UserID *string `json:"user_id"`

	Action string `json:"action"`

	Resource string `json:"resource"`

	ResourceID *string `json:"resource_id"`

	RequestID *string `json:"request_id"`

	IP *string `json:"ip"`

	UserAgent *string `json:"user_agent"`

	Result string `json:"result"`

	ErrorMessage *string `json:"error_message"`
}

func FromModelList(models []*model.AuditLog) []AuditLogResponse {
	var res []AuditLogResponse
	for _, model := range models {
		res = append(res, FromModel(model))
	}
	return res
}
func FromModel(model *model.AuditLog) AuditLogResponse {
	return AuditLogResponse{
		ID:           model.ID,
		UserID:       model.UserID,
		Action:       model.Action,
		Resource:     model.Resource,
		ResourceID:   model.ResourceID,
		RequestID:    model.RequestID,
		IP:           model.IP,
		UserAgent:    model.UserAgent,
		Result:       model.Result,
		ErrorMessage: model.ErrorMessage,
	}
}
