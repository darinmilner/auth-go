package repository

import (
	"auth/api/models"
	"context"
	"database/sql"

	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AllUsers returns all users
func (m *DBRepo) AllUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT id, name, email, username, updated_at FROM users
		where deleted_at is null`

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		s := &models.User{}
		err = rows.Scan(&s.ID, &s.Name, &s.Email, &s.UserName, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		// Append it to the slice
		users = append(users, s)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

// GetUserById returns a user by id
func (m *DBRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT id, name, email, username,
			created_at, updated_at
			FROM users where id = $1`
	row := m.DB.QueryRowContext(ctx, stmt, id)

	var u models.User

	err := row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.UserName,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		log.Println("DB can't get user info ", err)
		return u, err
	}

	return u, nil
}

// GetUserById returns a user by id
func (m *DBRepo) GetUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT id, name, email, username,
			created_at, updated_at, password_reset_code
			FROM users where email = $1`
	row := m.DB.QueryRowContext(ctx, stmt, email)

	var u models.User

	err := row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.UserName,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.PasswordResetCode,
	)

	if err != nil {
		log.Println("DB can't get user info ", err)
		return u, err
	}

	return u, nil
}

// Authenticate authenticates
func (m *DBRepo) Authenticate(username, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string
	//var userActive int

	query := `
		select 
		    id, password 
		from 
			users 
		where 
			username = $1
	and deleted_at is null`

	row := m.DB.QueryRowContext(ctx, query, username)
	err := row.Scan(&id, &hashedPassword)

	if err == sql.ErrNoRows {
		return 0, "", ErrInvalidCredentials
	} else if err != nil {
		log.Println(err)
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		log.Print(err)
		return 0, "", ErrInvalidCredentials
	} else if err != nil {
		log.Println(err)
		return 0, "", err
	}

	// Otherwise, the password is correct. Return the user ID and hashed password.
	return id, hashedPassword, nil
}

// Insert method to add a new record to the users table.
func (m *DBRepo) InsertUser(u models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	password := []byte(u.Password)

	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 12)
	if err != nil {
		log.Print(err)
		return 0, err
	}

	stmt := `
	INSERT INTO users 
	    (
		name,  
		email,
		username, 
		password, 
		created_at,
		updated_at
		)
    VALUES($1, $2, $3, $4, $5, $6) returning id `

	var newId int
	err = m.DB.QueryRowContext(ctx, stmt,
		u.Name,
		u.Email,
		u.UserName,
		hashedPassword,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		log.Print(err)
		return 0, err
	}

	return newId, err
}

// UpdateUser updates a user by id
func (m *DBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		update 
			users 
		set 
			name = $1, 
			email = $2, 
			username = $3
			updated_at = $4
		where
			id = $5`

	_, err := m.DB.ExecContext(ctx, stmt,
		u.Name,
		u.Email,
		u.UserName,
		time.Now(),
		u.ID,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// UpdateUserPassword updates a user by id
func (m *DBRepo) UpdateUserPassword(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		update 
			users 
		set 
			password = $1,
			updated_at = $2
		where
			email = $3 and password_reset_code = $4
			`

	_, err := m.DB.ExecContext(ctx, stmt,
		u.Password,
		time.Now(),
		u.Email,
		u.PasswordResetCode,
	)
	log.Println(u.Email)

	if err == sql.ErrNoRows {
		log.Println("User not in DB")
		return err
	}
	if err != nil {
		log.Println(err)
		return err
	}
	u.PasswordResetCode = ""
	newStmt := `
		update 
			users 
		set 
			password_reset_code = $1
		where
			email = $2 and password = $3
			`

	_, err = m.DB.ExecContext(ctx, newStmt,
		"",
		u.Email,
		u.Password,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// UpdateUser updates a user by id
func (m *DBRepo) AddResetPasswordCodeToUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		update 
			users 
		set 
			password_reset_code = $1
		where
			email = $2`

	_, err := m.DB.ExecContext(ctx, stmt,
		u.PasswordResetCode,
		u.Email,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// DeleteUser sets a user to deleted by populating deleted_at value
func (m *DBRepo) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update users set deleted_at = $1 where id = $2`

	_, err := m.DB.ExecContext(ctx, stmt, time.Now(), id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// ChangePassword resets a password
func (m *DBRepo) ChangePassword(id int, newPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		log.Println(err)
		return err
	}

	stmt := `update users set password = $1 where id = $2`
	_, err = m.DB.ExecContext(ctx, stmt, hashedPassword, id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
