package server

import (
	"context"
	"errors"
	"log"
	"payment-service/internal/service"
	"strings"

	structpb "github.com/golang/protobuf/ptypes/struct"
	pb "github.com/nexus-commerce/nexus-contracts-go/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedPaymentServiceServer
	Service *service.Service
}

func (s *Server) ProcessPayment(ctx context.Context, _ *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	userID, ok := ctx.Value("user-id").(int)
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "user id missing") // return FAILED_PRECONDITION status here as the system should never get into this state
	}

	userIDInt := int64(userID)

	pi, err := s.Service.ProcessPayment(ctx, userIDInt)
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
		Amount:       float32(pi.Amount) / 100,
		Currency:     strings.ToUpper(pi.Currency),
	}, nil
}

func (s *Server) GetPayments(ctx context.Context, _ *pb.GetPaymentsRequest) (*pb.GetPaymentsResponse, error) {
	userID, ok := ctx.Value("user-id").(int)
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "user id missing") // return FAILED_PRECONDITION status here as the system should never get into this state
	}

	userIDInt := int64(userID)

	transactions, err := s.Service.GetUserTransactions(ctx, userIDInt)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "error retrieving transactions")
	}

	var payments []*pb.Payment
	for _, t := range transactions {
		payment := &pb.Payment{
			Id:            t.ID,
			Amount:        float32(t.Amount),
			Currency:      string(t.Currency),
			Status:        string(t.Status),
			PaymentMethod: string(t.PaymentMethod),
		}

		if t.OrderID != nil {
			payment.OrderId = *t.OrderID
		}

		if t.GatewayTransactionID != nil {
			payment.GatewayTransactionId = &structpb.Value{
				Kind: &structpb.Value_StringValue{StringValue: *t.GatewayTransactionID},
			}
		}

		payments = append(payments, payment)
	}

	return &pb.GetPaymentsResponse{
		Payments: payments,
	}, nil
}
