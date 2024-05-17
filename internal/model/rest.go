package model

type Response struct {
	Data   interface{}           `json:"data"`
	Meta   Meta                  `json:"meta"`
	Errors []*ErrorFieldResponse `json:"errors,omitempty"`
}

type ResponseWithoutData struct {
	Meta   Meta                  `json:"meta"`
	Errors []*ErrorFieldResponse `json:"errors,omitempty"`
}

type Meta struct {
	Code    int    `json:"code" example:"3001"`
	Message string `json:"message" example:"success"`
}

type PaginationResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Meta
	Limit  int `json:"limit" example:"10"`
	Offset int `json:"offset" example:"0"`
	Total  int `json:"total" example:"100"`
}

type ErrorFieldResponse struct {
	Field      string `json:"field"`
	ErrMessage string `json:"message"`
	Tag        string `json:"tag"`
}

type ErrorResponse struct {
	Errors []ErrorFieldResponse `json:"errors"`
	Meta   Meta                 `json:"meta"`
}

type BusinessError struct {
	Code    int
	Message string
}

type UserIdentity struct {
	UserNik      string
	RequestID    string
	CacheControl string
}
