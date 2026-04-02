package request

type ProjectInviteRequest struct {
	ProjectID string `json:"project_id"`
	Email     string `json:"email"`
}
