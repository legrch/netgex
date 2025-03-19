package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/legrch/netgex/server"
)

func main() {
	// Create logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Create context that cancels on SIGINT or SIGTERM
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Example 1: Using environment variables for JSON configuration
	// Set environment variables before running:
	// export JSON_USE_PROTO_NAMES=true
	// export JSON_EMIT_UNPOPULATED=true
	// export JSON_USE_ENUM_NUMBERS=false
	// export JSON_ALLOW_PARTIAL=true
	// export JSON_MULTILINE=true
	// export JSON_INDENT="    "
	// app := server.NewServer(
	// 	server.WithLogger(logger),
	// )

	// Example 2: Using direct configuration
	// jsonConfig := runtime.ServeMuxOption{
	// 	UseProtoNames:   true,
	// 	EmitUnpopulated: true,
	// 	UseEnumNumbers:  false,
	// 	AllowPartial:    true,
	// 	Multiline:       true,
	// 	Indent:          "    ",
	// }

	app := server.NewServer(
		server.WithLogger(logger),
		server.WithGatewayMuxOptions(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
					UseEnumNumbers:  false,
					AllowPartial:    true,
					Multiline:       true,
					Indent:          "    ",
				},
			}),
		),
	)

	// Run the application
	if err := app.Run(ctx); err != nil {
		logger.Error("application error", "error", err)
		os.Exit(1)
	}
}
