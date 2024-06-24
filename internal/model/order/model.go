package order

type (
	OrderModel struct {
		ID         int64 `db:"id"`
		UserID     int64 `db:"user_id"`
		TotalQty   int64 `db:"total_quantity"`
		TotalPrice int64 `db:"total_price"`
	}

	OrderDetailModel struct {
		ID      int64 `db:"id"`
		OrderID int64 `db:"order_id"`
		BookID  int64 `db:"book_id"`
		Qty     int64 `db:"quantity"`
		Price   int64 `db:"price"`
	}
)

// Requests
type (
	CreateOrderRequest struct {
		Details []CreateOrderDetailRequest `json:"details"`
	}

	CreateOrderDetailRequest struct {
		BookID int64 `json:"book_id"`
		Qty    int64 `json:"quantity"`
		Price  int64 `json:"price"`
	}
)

// Responses
type (
	OrderResponse struct {
		ID         int64                 `json:"id"`
		UserID     int64                 `json:"user_id"`
		TotalQty   int64                 `json:"total_quantity"`
		TotalPrice int64                 `json:"total_price"`
		Details    []OrderDetailResponse `json:"details"`
	}

	OrderDetailResponse struct {
		ID     int64 `json:"id"`
		BookID int64 `json:"book_id"`
		Qty    int64 `json:"quantity"`
		Price  int64 `json:"price"`
	}
)
