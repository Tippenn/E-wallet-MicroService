package response

type CheckBalanceResponse struct {
	UserId  string `json:"userId" validate:"required"`
	Balance string `json:"balance"`
}
