package z

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type middlewaresRegistry struct{}

type responseWriter struct {
	http.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

var Middlewares = middlewaresRegistry{}

type LoggingConfig struct {
	LogRequestBody  bool
	LogResponseBody bool
	LogFilePath     string
}

func (mr middlewaresRegistry) Logging() MiddlewareFunc {
	return Middlewares.LoggingWithCfg(LoggingConfig{})
}

func (mr middlewaresRegistry) LoggingWithCfg(cfg LoggingConfig) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
			if cfg.LogFilePath != "" {
				f, err := openLogFile(cfg.LogFilePath)
				if err == nil {
					logger := slog.New(slog.NewJSONHandler(f, nil))
					slog.SetDefault(logger)
				} else {
					slog.Error("Failed to open log file", "err", err)
				}
			}

			start := time.Now()

			var requestBody []byte
			if cfg.LogRequestBody && z.r.Body != nil {
				var err error
				requestBody, err = io.ReadAll(z.r.Body)
				if err != nil {
					slog.Error("Error reading request body", "err", err)
				}
				z.r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			}

			writer := &responseWriter{
				ResponseWriter: z.rw, body: &bytes.Buffer{},
			}
			z.rw = writer

			next(z)

			latency := time.Since(start)
			reqID := z.r.Header.Get("X-Request-ID")

			logAttrs := []slog.Attr{
				slog.String("method", z.r.Method),
				slog.String("path", z.r.URL.Path),
				slog.Int("status", writer.status),
				slog.Duration("latency", latency),
			}

			if reqID != "" {
				logAttrs = append(logAttrs, slog.String("request_id", reqID))
			}

			if cfg.LogRequestBody && len(requestBody) > 0 {
				logAttrs = append(logAttrs, slog.String("request_body", string(requestBody)))
			}
			if cfg.LogResponseBody && writer.body.Len() > 0 {
				logAttrs = append(logAttrs, slog.String("response_body", writer.body.String()))
			}
			args := make([]any, len(logAttrs))
			for i, attr := range logAttrs {
				args[i] = attr
			}
			slog.Info("Request handled", args...)
		}
	}
}

func openLogFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}

type RecoveryConfig struct {
	LogPanic bool
}

func (middlewaresRegistry) Recovery() MiddlewareFunc {
	return Middlewares.RecoveryWithCfg(RecoveryConfig{LogPanic: true})
}

func (middlewaresRegistry) RecoveryWithCfg(cfg RecoveryConfig) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
			defer func() {
				if err := recover(); err != nil {
					if cfg.LogPanic {
						log.Printf("Recovered from panic: %v", err)
					}
					z.String(http.StatusInternalServerError, "Internal Server Error")
				}
			}()
			next(z)
		}
	}
}

type RequestIDConfig struct {
	HeaderName string
}

func (middlewaresRegistry) RequestID() MiddlewareFunc {
	return Middlewares.RequestIDWithCfg(RequestIDConfig{HeaderName: "X-Request-ID"})
}

func (middlewaresRegistry) RequestIDWithCfg(cfg RequestIDConfig) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
			reqID := z.r.Header.Get(cfg.HeaderName)
			if reqID == "" {
				reqID = generateRequestID()
			}
			z.r.Header.Set(cfg.HeaderName, reqID)
			z.rw.Header().Set(cfg.HeaderName, reqID)
			next(z)
		}
	}
}

func generateRequestID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return "req-" + hex.EncodeToString(b)
}

type CORSConfig struct {
	AllowOrigin      string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
	MaxAge           int
}

func (middlewaresRegistry) CORS() MiddlewareFunc {
	return Middlewares.CORSWithCfg(CORSConfig{
		AllowOrigin:      "*",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization",
		AllowCredentials: true,
		MaxAge:           3600,
	})
}

func (middlewaresRegistry) CORSWithCfg(cfg CORSConfig) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
			z.rw.Header().Set("Access-Control-Allow-Origin", cfg.AllowOrigin)
			z.rw.Header().Set("Access-Control-Allow-Methods", cfg.AllowMethods)
			z.rw.Header().Set("Access-Control-Allow-Headers", cfg.AllowHeaders)
			if cfg.AllowCredentials {
				z.rw.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if cfg.MaxAge > 0 {
				z.rw.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))
			}
			if z.r.Method == "OPTIONS" {
				z.rw.WriteHeader(http.StatusNoContent)
				return
			}
			next(z)
		}
	}
}

type SecurityHeadersConfig struct {
	ContentTypeOptions    string
	FrameOptions          string
	XSSProtection         string
	StrictTransportPolicy string
	ContentSecurityPolicy string
	ReferrerPolicy        string
}

func (middlewaresRegistry) SecurityHeaders() MiddlewareFunc {
	return Middlewares.SecurityHeadersWithCfg(SecurityHeadersConfig{
		ContentTypeOptions:    "nosniff",
		FrameOptions:          "DENY",
		XSSProtection:         "1; mode=block",
		StrictTransportPolicy: "max-age=31536000; includeSubDomains",
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "no-referrer",
	})
}

func (middlewaresRegistry) SecurityHeadersWithCfg(cfg SecurityHeadersConfig) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
			if cfg.ContentTypeOptions != "" {
				z.rw.Header().Set("X-Content-Type-Options", cfg.ContentTypeOptions)
			}
			if cfg.FrameOptions != "" {
				z.rw.Header().Set("X-Frame-Options", cfg.FrameOptions)
			}
			if cfg.XSSProtection != "" {
				z.rw.Header().Set("X-XSS-Protection", cfg.XSSProtection)
			}
			if cfg.StrictTransportPolicy != "" {
				z.rw.Header().Set("Strict-Transport-Security", cfg.StrictTransportPolicy)
			}
			if cfg.ContentSecurityPolicy != "" {
				z.rw.Header().Set("Content-Security-Policy", cfg.ContentSecurityPolicy)
			}
			if cfg.ReferrerPolicy != "" {
				z.rw.Header().Set("Referrer-Policy", cfg.ReferrerPolicy)
			}
			next(z)
		}
	}
}

type TimeoutConfig struct {
	Timeout time.Duration
}

func (middlewaresRegistry) Timeout() MiddlewareFunc {
	return Middlewares.TimeoutWithCfg(TimeoutConfig{Timeout: 30 * time.Second})
}

func (middlewaresRegistry) TimeoutWithCfg(cfg TimeoutConfig) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
			ctx, cancel := context.WithTimeout(z.r.Context(), cfg.Timeout)
			defer cancel()

			z.r = z.r.WithContext(ctx)

			done := make(chan struct{})
			go func() {
				next(z)
				close(done)
			}()

			select {
			case <-done:
			case <-ctx.Done():
				if ctx.Err() == context.DeadlineExceeded {
					z.String(http.StatusGatewayTimeout, "Request timed out")
				}
			}
		}
	}
}
