package rpc

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
)

type Grpc struct {
	Address string
}

func NewGrpc(address string) *Grpc {
	return &Grpc{Address: address}
}

func (g *Grpc) Register(opts []grpc.ServerOption, register func(s *grpc.Server)) error {
	lis, err := net.Listen("tcp", g.Address)
	if err != nil {
		return err
	}
	s := grpc.NewServer(opts...)
	register(s)
	fmt.Println("grpc connect success! listen address: " + g.Address)
	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
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

func Err(code codes.Code, msg string, stack []byte) error {
	s := status.New(code, msg)
	st, _ := s.WithDetails(&any.Any{
		Value:  stack,
	})
	return st.Err()
}
