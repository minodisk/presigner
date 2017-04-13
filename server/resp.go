package server

type Resp struct {
	code int         `json:"code"`
	Body interface{} `json:"body"`
}

func NewResp(code int, body interface{}) *Resp {
	return &Resp{
		code: code,
		Body: body,
	}
}

func NewErrorResp(code int, err error) *Resp {
	return &Resp{
		code: code,
		Body: Error{err.Error()},
	}
}

func (r *Resp) Code() int {
	return r.code
}

type Error struct {
	Error string `json:"error"`
}
