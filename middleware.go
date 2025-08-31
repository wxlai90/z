package z

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
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
}

func (mr middlewaresRegistry) Logging() MiddlewareFunc {
	return Middlewares.LoggingWithCfg(LoggingConfig{})
}

func (mr middlewaresRegistry) LoggingWithCfg(cfg LoggingConfig) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
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
