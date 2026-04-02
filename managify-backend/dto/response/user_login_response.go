package response

type UserLoginResponse struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}
