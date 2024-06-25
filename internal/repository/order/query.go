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
			od.id, od.order_id, od.book_id, od.quantity, od.price
		FROM
			order_details od
		JOIN
			orders o
		ON
			od.order_id = o.id
		WHERE
			od.order_id = ?
		AND
			o.user_id = ?
	`
)
