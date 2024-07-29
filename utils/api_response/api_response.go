package apiResponse

type Response struct {
	Message     string      `json:"message,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	Errors      interface{} `json:"errors,omitempty"`
	Error       string      `json:"error,omitempty"`
	CurrentPage int         `json:"currentPage,omitempty"`
	TotalPage   int         `json:"totalPage,omitempty"`
	Size        int         `json:"size,omitempty" validate:"min=1"`
}
