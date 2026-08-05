package utils

import (
	"context"
)

// RequestContext chứa tất cả thông tin của request
type RequestContext struct {
	Tid       string // Trace ID
	Sub       string // Subject/User ID
	DeviceID  string // Device ID
	IPAddress string // Client IP
	UserAgent string // User Agent
	SessionID string // Session ID
	Platform  string // Platform: web, ios, android, etc.
	Locale    string // Locale/Language (vd: vi-VN, en-US)
	Timezone  string // Timezone (vd: Asia/Ho_Chi_Minh)
}

// NewRequestContext tạo mới RequestContext
func NewRequestContext(keyRequest, sub, tid, deviceID, ipAddress, userAgent string) *RequestContext {
	return &RequestContext{
		Sub:       sub,
		Tid:       tid,
		DeviceID:  deviceID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
}

// NewRequestContextFull tạo mới RequestContext với đầy đủ thông tin
func NewRequestContextFull(
	sub, tid, deviceID, ipAddress, userAgent,
	shardID, sessionID, requestID, platform, locale, timezone string,
) *RequestContext {
	return &RequestContext{
		Sub:       sub,
		Tid:       tid,
		DeviceID:  deviceID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		SessionID: sessionID,
		Platform:  platform,
		Locale:    locale,
		Timezone:  timezone,
	}
}

// SaveRequestContext lưu RequestContext vào context
func SaveRequestContext(ctx context.Context, reqCtx *RequestContext, keyReq string) context.Context {
	return context.WithValue(ctx, keyReq, reqCtx)
}

// GetRequestContext lấy RequestContext từ context
func GetRequestContext(ctx context.Context, keyReq string) *RequestContext {
	val := ctx.Value(keyReq)
	if val == nil {
		return &RequestContext{}
	}
	r, ok := val.(*RequestContext)
	if !ok {
		return &RequestContext{}
	}
	return r
}

// GetSub lấy Subject từ context
func GetSub(ctx context.Context, keyReq string) string {
	reqCtx := GetRequestContext(ctx, keyReq)
	return reqCtx.Sub
}

// GetTid lấy Trace ID từ context
func GetTid(ctx context.Context, keyReq string) string {
	reqCtx := GetRequestContext(ctx, keyReq)
	return reqCtx.Tid
}

// GetDeviceID lấy Device ID từ context
func GetDeviceID(ctx context.Context, keyReq string) string {
	reqCtx := GetRequestContext(ctx, keyReq)
	return reqCtx.DeviceID
}

// GetSessionID lấy Session ID từ context
func GetSessionID(ctx context.Context, keyReq string) string {
	reqCtx := GetRequestContext(ctx, keyReq)
	return reqCtx.SessionID
}

// GetIPAddress lấy IP Address từ context
func GetIPAddress(ctx context.Context, keyReq string) string {
	reqCtx := GetRequestContext(ctx, keyReq)
	return reqCtx.IPAddress
}

// GetUserAgent lấy User Agent từ context
func GetUserAgent(ctx context.Context, keyReq string) string {
	reqCtx := GetRequestContext(ctx, keyReq)
	return reqCtx.UserAgent
}

// GetPlatform lấy Platform từ context
func GetPlatform(ctx context.Context, keyReq string) string {
	reqCtx := GetRequestContext(ctx, keyReq)
	return reqCtx.Platform
}

// GetLocale lấy Locale từ context
func GetLocale(ctx context.Context, keyReq string) string {
	reqCtx := GetRequestContext(ctx, keyReq)
	return reqCtx.Locale
}

// GetTimezone lấy Timezone từ context
func GetTimezone(ctx context.Context, keyReq string) string {
	reqCtx := GetRequestContext(ctx, keyReq)
	return reqCtx.Timezone
}
