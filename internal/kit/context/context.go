package context

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey int

const (
	ctxKeyUserID ctxKey = iota
	ctxKeyIP
	ctxKeyUserAgent
)

func PutUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxKeyUserID, userID)
}

func MustUserID(ctx context.Context) uuid.UUID {
	userID, ok := fromCtx[uuid.UUID](ctx, ctxKeyUserID)
	if !ok {
		panic("no user id in context")
	}

	return userID
}

func PutIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, ctxKeyIP, ip)
}

func IP(ctx context.Context) *string {
	ip, ok := fromCtx[string](ctx, ctxKeyIP)
	if !ok {
		return nil
	}

	return &ip
}

func PutUserAgent(ctx context.Context, userAgent string) context.Context {
	return context.WithValue(ctx, ctxKeyUserAgent, userAgent)
}

func UserAgent(ctx context.Context) *string {
	userAgent, ok := fromCtx[string](ctx, ctxKeyUserAgent)
	if !ok {
		return nil
	}

	return &userAgent
}

func fromCtx[T any](ctx context.Context, key ctxKey) (T, bool) {
	value := ctx.Value(key)
	if v, ok := value.(T); ok {
		return v, true
	}

	return *new(T), false
}
