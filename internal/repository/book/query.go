package book

var (
	queryGetAll = `
		SELECT 
			id, title, author, description, price
		FROM 
			books
	`
)
