package middleware

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/google/uuid"
)

const (
	// RequestIDHeader 请求ID的HTTP头
	RequestIDHeader = "X-Request-ID"
)

type requestIDKey struct{}

// RequestID 请求ID中间件
func RequestID() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var requestID string

			// 从传输层获取请求ID
			if tr, ok := transport.FromServerContext(ctx); ok {
				requestID = getRequestIDFromHeader(tr)
			}

			// 生成新的请求ID
			if requestID == "" {
				requestID = generateUniqueRequestID()
			}

			// 清理请求ID（移除空格等）
			requestID = strings.TrimSpace(requestID)

			// 存储到context中
			ctx = context.WithValue(ctx, requestIDKey{}, requestID)

			// 设置到响应头中
			if tr, ok := transport.FromServerContext(ctx); ok {
				setResponseHeader(tr, requestID)
			}

			return handler(ctx, req)
		}
	}
}

// 从头部获取请求ID
func getRequestIDFromHeader(tr transport.Transporter) string {
	headers := tr.RequestHeader()
	if headers == nil {
		return ""
	}

	// 使用正确的方法获取头部值
	headerValue := headers.Get(RequestIDHeader)
	if headerValue == "" {
		return ""
	}

	return headerValue
}

// 设置响应头
func setResponseHeader(tr transport.Transporter, requestID string) {
	headers := tr.ReplyHeader()
	if headers == nil {
		return
	}

	headers.Set(RequestIDHeader, requestID)
}

// 生成唯一请求ID
func generateUniqueRequestID() string {
	return uuid.New().String()
}

// FromContext 从context中获取请求ID
func FromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return generateUniqueRequestID()
}

// WithContext 将请求ID设置到context中
func WithContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}
