package usecase

import (
	"github.com/gilbertom/desafio-clean-architecture/internal/entity"
)

type QueryOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewQueryOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *QueryOrderUseCase {
	return &QueryOrderUseCase{
		OrderRepository: OrderRepository,
	}
}

// FindAll retrieves all orders.
func (c *QueryOrderUseCase) FindAll() ([]OrderOutputDTO, error) {
	orders, err := c.OrderRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var ordersOutput []OrderOutputDTO
	for _, order := range orders {
		ordersOutput = append(ordersOutput, OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.Price + order.Tax,
		})
	}
	
	return ordersOutput, nil
}
