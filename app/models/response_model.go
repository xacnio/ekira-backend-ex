package models

import "ekira-backend/app/errs"

type ResponseOK struct {
	ErrorType error                  `json:"-"`
	Error     interface{}            `json:"error" extensions:"x-nullable,x-example=null"`
	Headers   map[string]interface{} `json:"headers" swaggertype:"object"`
	IsSuccess bool                   `json:"isSuccess" example:"true"`
	Result    interface{}            `json:"result"`
}

type ResponseErr struct {
	ErrorType error                  `json:"-"`
	Error     interface{}            `json:"error" swaggertype:"string" example:"error message"`
	Headers   map[string]interface{} `json:"headers" swaggertype:"object"`
	IsSuccess bool                   `json:"isSuccess" example:"false"`
	Result    interface{}            `json:"result" extensions:"x-nullable,x-example=null"`
}

func NewResponseString(result string) *ResponseOK {
	response := &ResponseOK{
		ErrorType: nil,
		Error:     nil,
		Headers:   map[string]interface{}{},
		IsSuccess: true,
		Result:    &result,
	}
	return response
}

func NewResponseOK(result any) *ResponseOK {
	response := &ResponseOK{
		ErrorType: nil,
		Error:     nil,
		Headers:   map[string]interface{}{},
		IsSuccess: true,
		Result:    result,
	}
	return response
}

func NewResponseError(error *errs.Error) *ResponseErr {
	response := &ResponseErr{
		ErrorType: error.Error,
		Error:     error.Description,
		Headers:   map[string]interface{}{},
		IsSuccess: false,
		Result:    nil,
	}
	return response
}

func NewResponseErr(error error) *ResponseErr {
	response := &ResponseErr{
		ErrorType: error,
		Error:     error.Error(),
		Headers:   map[string]interface{}{},
		IsSuccess: false,
		Result:    nil,
	}
	return response
}

func (r *ResponseErr) SetHeader(key string, data interface{}) *ResponseErr {
	r.Headers[key] = data
	return r
}

func (r *ResponseErr) ToInterface() interface{} {
	return r
}
