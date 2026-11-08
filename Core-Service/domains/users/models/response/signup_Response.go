package response

type SignUpResponse struct {
	UserId  string `json:"userID"`
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required"`
	Message string `json:"message" validate:"required"`
}
