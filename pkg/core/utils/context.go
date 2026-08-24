package utils

import (
	"context"
	wcmd_utils "github.com/DVV-15324/witches/cmd/utils"
)

type RequestContext struct {
	Tid string
	Sub string
}

func NewRequestContext(sub, tid string) *RequestContext {
	return &RequestContext{
		Sub: sub,
		Tid: tid,
	}
}

func NewRequestContextFull(sub, tid string) *RequestContext {
	return &RequestContext{
		Sub: sub,
		Tid: tid,
	}
}

func SaveRequestContext(ctx context.Context, reqCtx *RequestContext, cfg *wcmd_utils.Config) context.Context {
	return context.WithValue(ctx, cfg.RequestKey, reqCtx)
}

func GetRequestContext(ctx context.Context, cfg *wcmd_utils.Config) *RequestContext {
	val := ctx.Value(cfg.RequestKey)
	if val == nil {
		return &RequestContext{}
	}
	r, ok := val.(*RequestContext)
	if !ok {
		return &RequestContext{}
	}
	return r
}

func GetSub(ctx context.Context, cfg *wcmd_utils.Config) string {
	reqCtx := GetRequestContext(ctx, cfg)
	return reqCtx.Sub
}

func GetTid(ctx context.Context, cfg *wcmd_utils.Config) string {
	reqCtx := GetRequestContext(ctx, cfg)
	return reqCtx.Tid
}
