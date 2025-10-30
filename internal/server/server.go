package server

import (
	"context"
	"errors"
	pb "github.com/nexus-commerce/nexus-contracts-go/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"payment-service/internal/service"
)

type Server struct {
	pb.UnimplementedPaymentServiceServer
	Service *service.Service
}

func (s *Server) ProcessPayment(ctx context.Context, _ *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	pi, err := s.Service.ProcessPayment(ctx)
	if err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, service.ErrEmptyCart):
			return nil, status.Error(codes.FailedPrecondition, "cart is empty")
		case errors.Is(err, service.ErrProcessingPayment):
			return nil, status.Error(codes.Internal, "error processing pb")
		}
		return nil, status.Error(codes.Internal, "unknown error processing pb")
	}

	return &pb.ProcessPaymentResponse{
		ClientSecret: pi.ClientSecret,
		Amount:       float32(pi.Amount),
		Currency:     pi.Currency,
	}, nil
}
