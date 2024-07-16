package database

import (
	"database/sql"

	"github.com/gilbertom/desafio-clean-architecture/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.Db.QueryRow("Select count(*) from orders").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// FindAll retrieves all orders from the database.
func (r *OrderRepository) FindAll() ([]entity.Order, error) {
	var orders []entity.Order
	query := "SELECT id, price, tax, final_price FROM orders"
	rows, err := r.Db.Query(query)
	if err != nil {
		return orders, err
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		err := rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return orders, err
	}

	return orders, nil
}