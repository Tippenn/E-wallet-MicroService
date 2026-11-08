package request

type LogoutRequest struct {
	UserId string `json:"user_Id" validate:"required"`
}
