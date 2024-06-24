package order

var (
	queryCreate = `
		INSERT INTO orders 
			(user_id, total_quantity, total_price)
		VALUES
			(?, ?, ?)
		RETURNING
			id
	`

	queryCreateDetail = `
		INSERT INTO order_details
			(order_id, book_id, quantity, price)
		VALUES
			%s
		RETURNING
			id
	`

	queryGetAllOrder = `
		SELECT
			id, user_id, total_quantity, total_price
		FROM
			orders
		WHERE
			user_id = ?
	`

	queryGetOrderDetail = `
		SELECT
			id, order_id, book_id, quantity, price
		FROM
			order_details
		WHERE
			order_id = ?
	`
)
