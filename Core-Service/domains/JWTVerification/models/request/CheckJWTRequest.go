package request

type CheckJWTRequest struct {
	UserId string `json:"user_Id" validate:"required"`
}
