package usecase

import (
	"challenge-arch/internal/entity"
	"challenge-arch/pkg/events"
)

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	listOrders      events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewListOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	listOrders events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
		listOrders:      listOrders,
		EventDispatcher: EventDispatcher,
	}
}

func (c *ListOrdersUseCase) Execute() ([]OrderOutputDTO, error) {

	list, err := c.OrderRepository.List()
	if err != nil {
		return nil, err
	}
	listDtos := make([]OrderOutputDTO, len(list))

	for idx, order := range list {
		listDtos[idx] = OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
	}

	c.listOrders.SetPayload(listDtos)
	c.EventDispatcher.Dispatch(c.listOrders)

	return listDtos, nil
}
