package rpc

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
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
	s := status.Convert(err)
	var stack string
	for _, detail := range s.Details() {
		switch t := detail.(type) {
		case *any.Any:
			stack = string(t.Value)
		}
	}

	errReport(md, req, stack, s)
	return status.Error(s.Code(), s.Message())
}

func Err(code codes.Code, msg string) error {
	s := status.New(code, msg)
	st, _ := s.WithDetails(&any.Any{
		Value: stack(msg),
	})
	return st.Err()
}

func stack(msg string) []byte {
	stack := fmt.Sprintf("%+v\n", errors.New(msg))
	fmt.Println(stack)
	return []byte(stack)
}
