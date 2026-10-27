package utils

import (
	"context"
	wcmd_utils "github.com/DVV-15324/witches/cmd/utils"
)

type RequestContext struct {
	Tid   string
	Sub   string
	RQCxt interface{}
}

func NewRequestContext(sub, tid string) *RequestContext {
	return &RequestContext{
		Sub: sub,
		Tid: tid,
	}
}

func (r *RequestContext) SetExtra(v interface{}) {
	r.RQCxt = v
}

func GetExtra[T any](ctx context.Context, cfg *wcmd_utils.Config) (T, bool) {
	var zero T
	r := GetRequestContext(ctx, cfg)
	if r.RQCxt == nil {
		return zero, false
	}
	v, ok := r.RQCxt.(T)
	return v, ok
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
