package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

var db *sql.DB

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	User User
}

type User struct {
	ID                int       `json:"id"`
	Email             string    `json:"email"`
	FirstName         string    `json:"first_name,omitempty"`
	LastName          string    `json:"last_name,omitempty"`
	Password          string    `json:"password"`
	Active            int       `json:"active"`
	RollId            int       `json:"roll_id"`
	Roll              string    `json:"roll"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

// GetAll returns a slice of all users, sorted by last name
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select users.id, users.email, users.first_name, users.last_name, users.hashad_password, users.active, users.roll_id, users.password_changed_at, users.updated_at, users.created_at,roll.roll from users
	INNER JOIN roll on users.roll_id = roll.id`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.RollId,
			&user.Roll,
			&user.PasswordChangedAt,
			&user.UpdatedAt,
			&user.CreatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// GetByEmail returns one user by email
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select users.id, users.email, users.first_name, users.last_name, users.hashad_password, users.active, users.roll_id,roll.roll, users.password_changed_at, users.updated_at, users.created_at from users
	INNER JOIN roll on users.roll_id = roll.id
	WHERE users.email = $1`

	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.RollId,
		&user.Roll,
		&user.PasswordChangedAt,
		&user.UpdatedAt,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetOne returns one user by id
func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select users.id, users.email, users.first_name, users.last_name, users.active, users.roll_id, roll.roll, users.password_changed_at, users.updated_at, users.created_at from users
	INNER JOIN roll on users.roll_id = roll.id
	where users.id = $1`

	var user User
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Active,
		&user.RollId,
		&user.Roll,
		&user.PasswordChangedAt,
		&user.UpdatedAt,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `UPDATE users SET
		email = $1,
		first_name = $2,
		last_name = $3,
		active = $4,
		roll_id = $5,
		updated_at = $6
		WHERE id = $7;
	`

	result, err := db.ExecContext(ctx, stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Active,
		u.RollId,
		time.Now(),
		u.ID,
	)

	if err != nil {
		log.Println(err)
		return err
	}
	inserted, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return errors.New("failed to update user")
	} else if inserted == 0 {
		fmt.Println("Updated value ", inserted)
		return errors.New("failed to update user")
	}

	return nil
}

// Delete deletes one user from the database, by User.ID
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID deletes one user from the database, by ID
func (u *User) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *User) Insert(user User) (int, error) {

	//ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	//defer cancel()
	log.Println("register plain password:", user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	log.Println("register hash password:", hashedPassword)
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `INSERT INTO users (email, first_name, last_name, hashad_password, active, roll_id) values ($1, $2, $3, $4, $5, $6) returning id`

	err = db.QueryRow(stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		user.RollId,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set password = $1, password_changed_at= $2 where id = $3`
	_, err = db.ExecContext(ctx, stmt, hashedPassword, time.Now(), u.ID)
	if err != nil {
		return err
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *User) PasswordMatches(plainText string) (bool, error) {
	log.Println(u.Password)
	log.Println(plainText)
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
