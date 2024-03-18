package responses

const (
	SUCCESS = "000"
)

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Ok(message string, data interface{}) *Response {
	return &Response{
		Code:    SUCCESS,
		Message: message,
		Data:    data,
	}
}
