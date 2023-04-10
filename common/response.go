package common

type Response struct {
	Code     int    `json:"code"`
	ErrorMsg string `json:"error_msg"`
	Data     any    `json:"data"`
}

func (r Response) Error() string {
	return r.ErrorMsg
}
