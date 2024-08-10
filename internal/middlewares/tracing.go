package middlewares

import (
	otel "github.com/hazkall/shadowlog/pkg/telemetry"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

func Tracing(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, span := otel.Tracer.Start(c.Request().Context(),
			c.Request().URL.Path,
			trace.WithAttributes(
				semconv.HTTPMethodKey.String(c.Request().Method),
				semconv.HTTPURLKey.String(c.Request().URL.Path),
			),
		)
		defer span.End()

		err := next(c)

		statusCode := c.Response().Status
		if statusCode >= 200 && statusCode < 400 {
			span.SetStatus(codes.Ok, "HTTP Success")
		} else {
			span.SetStatus(codes.Error, "HTTP Error")
		}

		return err
	}
}
