// Package greeterv1 contains mock implementations of the greeter service proto definitions
// In a real application, this would be generated from the proto file
package greeterv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SayHelloRequest is the request for SayHello
type SayHelloRequest struct {
	Name string `json:"name"`
}

// SayHelloResponse is the response from SayHello
type SayHelloResponse struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// GreeterServiceServer is the server API for GreeterService service.
type GreeterServiceServer interface {
	// SayHello sends a greeting to the requested name
	SayHello(context.Context, *SayHelloRequest) (*SayHelloResponse, error)
	// SayHelloStream streams greetings to the requested name
	SayHelloStream(*SayHelloRequest, GreeterService_SayHelloStreamServer) error
}

// UnimplementedGreeterServiceServer can be embedded to have forward compatible implementations.
type UnimplementedGreeterServiceServer struct{}

func (UnimplementedGreeterServiceServer) SayHello(context.Context, *SayHelloRequest) (*SayHelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}

func (UnimplementedGreeterServiceServer) SayHelloStream(*SayHelloRequest, GreeterService_SayHelloStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method SayHelloStream not implemented")
}

// GreeterService_SayHelloStreamServer is the server API for GreeterService_SayHelloStream service.
type GreeterService_SayHelloStreamServer interface {
	Send(*SayHelloResponse) error
	grpc.ServerStream
}

// RegisterGreeterServiceServer registers the server with the given gRPC server.
func RegisterGreeterServiceServer(s *grpc.Server, srv GreeterServiceServer) {
	s.RegisterService(&_GreeterService_serviceDesc, srv)
}

// RegisterGreeterServiceHandlerFromEndpoint registers the http handlers for service GreeterService to "mux".
func RegisterGreeterServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	// Register the HTTP handlers for path parameter style
	err := mux.HandlePath("GET", "/v1/greeter/{name}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		name := pathParams["name"]
		if name == "" {
			http.Error(w, "name parameter is required", http.StatusBadRequest)
			return
		}

		// Create a response directly (in a real implementation, this would call the gRPC service)
		resp := &SayHelloResponse{
			Message:   fmt.Sprintf("Hello, %s!", name),
			Timestamp: time.Now().Unix(),
		}

		// Write the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	if err != nil {
		return err
	}

	// Also register a handler for query parameter style
	return mux.HandlePath("GET", "/v1/greeter/hello", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "World" // Default name
		}

		// Create a response directly
		resp := &SayHelloResponse{
			Message:   fmt.Sprintf("Hello, %s!", name),
			Timestamp: time.Now().Unix(),
		}

		// Write the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}

var _GreeterService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "greeter.v1.GreeterService",
	HandlerType: (*GreeterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _GreeterService_SayHello_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SayHelloStream",
			Handler:       _GreeterService_SayHelloStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "greeter.proto",
}

func _GreeterService_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SayHelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServiceServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/greeter.v1.GreeterService/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServiceServer).SayHello(ctx, req.(*SayHelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GreeterService_SayHelloStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	in := new(SayHelloRequest)
	if err := stream.RecvMsg(in); err != nil {
		return err
	}
	return srv.(GreeterServiceServer).SayHelloStream(in, &greeterServiceSayHelloStreamServer{stream})
}

type greeterServiceSayHelloStreamServer struct {
	grpc.ServerStream
}

func (x *greeterServiceSayHelloStreamServer) Send(m *SayHelloResponse) error {
	return x.ServerStream.SendMsg(m)
}
