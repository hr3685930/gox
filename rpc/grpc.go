package rpc

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"runtime/debug"
)

type Grpc struct {
	G     *grpc.Server
	Debug bool
}

func NewGrpc(debug bool) *Grpc {
	return &Grpc{Debug: debug}
}

func (g *Grpc) Register(opts []grpc.ServerOption, register func(s *grpc.Server)) error {
	s := grpc.NewServer(opts...)
	register(s)
	g.G = s
	return nil
}

func (g *Grpc) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Println("grpc connect success! listen address: " + addr)
	return g.G.Serve(lis)
}

type GrpcError struct {
	Err   error
	Stack []byte `json:"-"`
}

func (h *GrpcError) Error() string {
	return h.Err.Error()
}

func (h *GrpcError) GetStack() string {
	return string(h.Stack)
}


type ErrorReport func(md metadata.MD, req interface{}, stack string, resp *status.Status)

func CustomErrInterceptor(errReport ErrorReport) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			md, _ := metadata.FromIncomingContext(ctx)
			err = ErrorHandler(md, req, err, errReport)
		}
		return
	}
}

func ErrorHandler(md metadata.MD, req interface{}, err error, errReport ErrorReport) error {
	var stack string
	if e, ok := err.(*GrpcError); ok {
		stack = string(e.Stack)
		err = e.Err
	} else {
		stack = string(debug.Stack())
	}

	s := status.Convert(err)
	errReport(md, req, stack, s)
	return status.Error(s.Code(), s.Message())
}

func Err(code codes.Code, msg string) *GrpcError {
	return &GrpcError{
		Err:   status.New(code, msg).Err(),
		Stack: []byte(fmt.Sprintf("%+v\n", errors.New(msg))),
	}
}
