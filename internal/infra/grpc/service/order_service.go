package service

import (
	"context"

	"github.com/gilbertom/desafio-clean-architecture/internal/infra/grpc/pb"
	"github.com/gilbertom/desafio-clean-architecture/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	QueryOrderUseCase  usecase.QueryOrderUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase,	queryOrderUseCase usecase.QueryOrderUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		QueryOrderUseCase: queryOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

// ListOrder retrieves a list of orders.
func (s *OrderService) ListOrder(ctx context.Context, in *pb.Blank) (*pb.OrderList, error) {
	orders, err := s.QueryOrderUseCase.FindAll()
	if err != nil {
		return nil, err
	}
	var ordersResponse []*pb.ListOrderResponse

	for _, order := range orders {
		ordersResponse = append(ordersResponse, &pb.ListOrderResponse{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		})
	}
	
	return &pb.OrderList{
		Orders: ordersResponse,
	}, nil
}