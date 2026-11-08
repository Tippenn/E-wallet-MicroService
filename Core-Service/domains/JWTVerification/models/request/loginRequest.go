package request

type LoginRequest struct {
	UserId      string `json:"user_Id" validate:"required"`
	AccessToken string `json:"access_Token" validate:"required"`
}
