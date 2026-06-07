package utils

import "context"

type RequestResponse struct {
	Sub string
	Tid string
}

func NewRequestResponse(sub string, tid string) *RequestResponse {
	return &RequestResponse{Sub: sub, Tid: tid}
}

type keyRequest string

var KeyReq keyRequest

func SaveRequestContext(cxt context.Context, r RequestResponse) context.Context {
	return context.WithValue(cxt, KeyReq, r)
}

func GetRequestContext(cxt context.Context) RequestResponse {
	r, _ := cxt.Value(KeyReq).(RequestResponse)
	return r
}
