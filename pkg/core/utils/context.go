package utils

import (
	"context"
	wcmd_utils "github.com/DVV-15324/witches/cmd/utils"
)

type RequestContext struct {
	Tid      string
	Sub      string
	DeviceID string
	Locale   string
	Timezone string
}

func NewRequestContext(sub, tid, deviceID string) *RequestContext {
	return &RequestContext{
		Sub:      sub,
		Tid:      tid,
		DeviceID: deviceID,
	}
}

func NewRequestContextFull(sub, tid, deviceID, locale, timezone string) *RequestContext {
	return &RequestContext{
		Sub:      sub,
		Tid:      tid,
		DeviceID: deviceID,
		Locale:   locale,
		Timezone: timezone,
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

func GetDeviceID(ctx context.Context, cfg *wcmd_utils.Config) string {
	reqCtx := GetRequestContext(ctx, cfg)
	return reqCtx.DeviceID
}

func GetLocale(ctx context.Context, cfg *wcmd_utils.Config) string {
	reqCtx := GetRequestContext(ctx, cfg)
	return reqCtx.Locale
}

func GetTimezone(ctx context.Context, cfg *wcmd_utils.Config) string {
	reqCtx := GetRequestContext(ctx, cfg)
	return reqCtx.Timezone
}
