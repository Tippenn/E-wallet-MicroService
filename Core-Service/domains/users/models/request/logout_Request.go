package request

type LogoutRequest struct {
	AccessToken string `json:"accessToken" validate:"required"`
}
