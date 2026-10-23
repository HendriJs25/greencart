package response

type Response struct {
	Status  string `json:"status"`
	Message any    `json:"message,omitempty"`
	Data    any    `json:"data"`
}
