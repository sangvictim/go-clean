package apiResponse

type Response struct {
	Message     string `json:"message,omitempty"`
	Data        any    `json:"data,omitempty"`
	Errors      any    `json:"errors,omitempty"`
	Error       string `json:"error,omitempty"`
	CurrentPage int    `json:"currentPage,omitempty"`
	TotalPage   int    `json:"totalPage,omitempty"`
	Limit       int    `json:"limit,omitempty" validate:"min=1"`
}
