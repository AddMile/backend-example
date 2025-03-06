package codegen

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"

	contextkit "github.com/AddMile/backend/internal/kit/context"
	ratelimitkit "github.com/AddMile/backend/internal/kit/ratelimit"
)

type (
	handlerFunc    = nethttp.StrictHTTPHandlerFunc
	middlewareFunc = nethttp.StrictHTTPMiddlewareFunc
)

func Recover(f handlerFunc, operationID string) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic occurred: %v\n%s", err, debug.Stack())

				w.Header().Set("Connection", "close")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		return f(ctx, w, r, request)
	}
}

func CORS(origins string) middlewareFunc {
	allowedOrigins := make(map[string]struct{})
	for _, o := range strings.Split(origins, ";") {
		allowedOrigins[o] = struct{}{}
	}

	return func(f handlerFunc, operationID string) handlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
			origin := r.Header.Get("Origin")

			_, ok := allowedOrigins[origin]
			if ok || origins == "*" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "HEAD, GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cache-Control, User-Agent")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "86400") // cache preflight requests for 1 day
			}

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				log.Println("CORS Middleware: Handled OPTIONS request")
				w.WriteHeader(http.StatusNoContent)

				return nil, nil
			}

			return f(ctx, w, r, request)
		}
	}
}

func WithUserAgent(f handlerFunc, operationID string) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		userAgent := r.Header.Get("User-Agent")

		ctx = contextkit.PutUserAgent(ctx, userAgent)

		return f(ctx, w, r, request)
	}
}

func WithIP(f handlerFunc, operationID string) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		ip := extractIP(r)

		ctx = contextkit.PutIP(ctx, ip)

		return f(ctx, w, r, request)
	}
}

func extractIP(r *http.Request) string {
	headers := []string{
		r.Header.Get("X-Forwarded-For"), // May contain a list of IPs
		r.Header.Get("X-Real-IP"),
		r.Header.Get("X-Client-IP"),
	}

	for _, h := range headers {
		if h != "" {
			ips := strings.Split(h, ",")
			ip := strings.TrimSpace(ips[0])

			if isValidIP(ip) {
				return ip
			}
		}
	}

	// Fallback to RemoteAddr if no headers contain a valid IP
	if r.RemoteAddr != "" {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil && isValidIP(ip) {
			return ip
		}
	}

	return ""
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func RateLimit(rl *ratelimitkit.RateLimiter) middlewareFunc {
	return func(f handlerFunc, operationID string) handlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				return nil, fmt.Errorf("cannot parse remote address: %w", err)
			}

			if !rl.AllowedByIP(ip) {
				w.WriteHeader(http.StatusTooManyRequests)

				return nil, nil
			}

			return f(ctx, w, r, request)
		}
	}
}
