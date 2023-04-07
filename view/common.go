package view

type Resp[Data any] struct {
	Data Data   `json:"data"`
	Code string `json:"code"`
	Msg  string `json:"msg"`
}
