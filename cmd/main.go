package main

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/hazkall/shadowlog/internal/request"
	"github.com/hazkall/shadowlog/internal/router"
	"github.com/hazkall/shadowlog/pkg/logger"
	otel "github.com/hazkall/shadowlog/pkg/telemetry"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/codes"
)

var (
	interval int
	port     int
)

func main() {
	var cmdShadowLog = &cobra.Command{
		Use:   "run --interval 2 --port 3000",
		Short: "Run shadowlog",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			logger.Start()

			slog.Info("Starting OpenTelemetry Tracing")
			t := otel.TraceInit(ctx, "")
			defer func() {
				if err := t.Shutdown(ctx); err != nil {
					panic(err)
				}
			}()

			slog.Info("Starting OpenTelemetry Metrics")
			m, exp := otel.MetricInit(ctx, "")

			defer func() {
				if err := exp.Shutdown(ctx); err != nil {
					panic(err)
				}
			}()

			slog.Info("Starting OpenTelemetry Go Runtime Metrics")
			otel.RuntimeStart(m)

			go RunServer()
			RunRequests(ctx)

			return nil
		},
	}

	cmdShadowLog.Flags().IntVarP(&interval, "interval", "i", 2, "Interval to make requests to the server")
	cmdShadowLog.Flags().IntVarP(&port, "port", "p", 3000, "Port to run the server on")
	cmdShadowLog.MarkFlagRequired("interval")
	cmdShadowLog.MarkFlagRequired("port")

	var rootCmd = &cobra.Command{Use: "shadowlog"}

	rootCmd.AddCommand(cmdShadowLog)
	rootCmd.Execute()
}

func RunServer() {
	r := router.NewRouter(port)
	r.Start()
}

func RunRequests(ctx context.Context) {
	h := map[string]string{
		"Content-Type": "application/json",
	}

	for {
		_, span := otel.Tracer.Start(ctx, "RunRequests")

		respCode, body, err := request.MakeHTTPGet("http://localhost:"+strconv.Itoa(port)+"/api/v1/deploy/shoulddeploy", h)
		if err != nil {
			slog.Error(err.Error(), slog.String("trace_id", span.SpanContext().TraceID().String()))
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}

		resp := fmt.Sprintf("Response: %d, Body: %s", respCode, string(body))
		slog.Info(resp, slog.String("trace_id", span.SpanContext().TraceID().String()))
		span.SetStatus(codes.Ok, resp)

		time.Sleep(time.Duration(interval) * time.Second)
		span.End()
	}
}
