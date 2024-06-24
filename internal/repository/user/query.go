package user

var (
	queryCreate = `
		INSERT INTO users 
			(email, password, name) 
		VALUES 
			(?, ?, ?) 
		RETURNING 
			id
	`

	queryGetByEmail = `
		SELECT 
			id, email, password
		FROM 
			users
		WHERE 
			email = ?
	`
)
