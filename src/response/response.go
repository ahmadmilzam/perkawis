package response

type (
	SuccessResponse struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	ErrorResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)
