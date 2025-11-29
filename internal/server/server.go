package server

import (
	"context"
	"errors"
	structpb "github.com/golang/protobuf/ptypes/struct"
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
	userID, ok := ctx.Value("user-id").(int)
	if !ok {
		return nil, status.Error(codes.Internal, "user id missing") // return FAILED_PRECONDITION status here as the system should never get into this state
	}

	userIDInt := int64(userID)

	t, pi, err := s.Service.ProcessPayment(ctx, userIDInt)
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
		Payment: &pb.Payment{
			Id:            t.ID,
			Amount:        float32(t.Amount),
			Currency:      pb.Currency(pb.Currency_value[string(t.Currency)]),
			Status:        pb.Status(pb.Status_value[string(t.Status)]),
			PaymentMethod: pb.PaymentMethod(pb.PaymentMethod_value[string(t.PaymentMethod)]),
			OrderId:       0,
			GatewayTransactionId: &structpb.Value{
				Kind: &structpb.Value_StringValue{StringValue: *t.GatewayTransactionID},
			},
		},
	}, nil
}

func (s *Server) GetPayments(ctx context.Context, _ *pb.GetPaymentsRequest) (*pb.GetPaymentsResponse, error) {
	userID, ok := ctx.Value("user-id").(int)
	if !ok {
		return nil, status.Error(codes.Internal, "user id missing")
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
			Currency:      pb.Currency(pb.Currency_value[string(t.Currency)]),
			Status:        pb.Status(pb.Status_value[string(t.Status)]),
			PaymentMethod: pb.PaymentMethod(pb.PaymentMethod_value[string(t.PaymentMethod)]),
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
