package models

import (
	"authentication/data/proto"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const dbTimeout = time.Second * 3

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

// UserService represents a PostgreSQL implementation of myapp.UserService.
type UserServiceStruct struct {
	DB *sql.DB
}

type UserService interface {
	GetOne(id int) (*User, error)
	GetAll() ([]User, error)
	GetByEmail(email string) (*User, error)
	Update(usr *User) error
	ValidateEmail(usr *User) error
	Delete(usr *User) error
	DeleteByID(id int) error
	Insert(user User) (User, error)
	ResetPassword(password string, usr *User) error
	PasswordMatches(plainText string, usr *User) (bool, error)
	ResetUsers() error
	ConvertToProtoUser(user User) proto.User
}

// GetAll returns a slice of all users, sorted by last name
func (us UserServiceStruct) GetAll() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select users.id, users.email, users.first_name, users.last_name, users.hashad_password, users.active, users.roll_id, users.password_changed_at, users.updated_at, users.created_at,roll.roll from users
	INNER JOIN roll on users.roll_id = roll.id`

	rows, err := us.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

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

		users = append(users, user)
	}

	return users, nil
}

// GetByEmail returns one user by email
func (us UserServiceStruct) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select users.id, users.email, users.first_name, users.last_name, users.hashad_password, users.active, users.roll_id,roll.roll, users.password_changed_at, users.updated_at, users.created_at from users
	INNER JOIN roll on users.roll_id = roll.id
	WHERE users.email = $1`

	var user User
	row := us.DB.QueryRowContext(ctx, query, email)

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
func (us UserServiceStruct) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select users.id, users.email, users.first_name, users.last_name, users.active, users.roll_id, roll.roll, users.password_changed_at, users.updated_at, users.created_at from users
	INNER JOIN roll on users.roll_id = roll.id
	where users.id = $1`

	var user User
	row := us.DB.QueryRowContext(ctx, query, id)

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
func (us UserServiceStruct) Update(usr *User) error {
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

	result, err := us.DB.ExecContext(ctx, stmt,
		usr.Email,
		usr.FirstName,
		usr.LastName,
		usr.Active,
		usr.RollId,
		time.Now(),
		usr.ID,
	)

	if err != nil {
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

// activate account if hash is correct
func (us UserServiceStruct) ValidateEmail(usr *User) error {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `UPDATE users SET
		active = $1,
		updated_at = $2
		WHERE email = $3;
	`

	result, err := us.DB.ExecContext(ctx, stmt,
		1,
		time.Now(),
		usr.Email,
	)

	if err != nil {
		log.Println(err)
		return err
	}
	inserted, err := result.RowsAffected()
	if err != nil {
		//log.Println(err)
		return errors.New("failed to update user")
	} else if inserted == 0 {
		//fmt.Println("Updated value ", inserted)
		return errors.New("failed to update user")
	}

	return nil
}

// Delete deletes one user from the database, by User.ID
func (us UserServiceStruct) Delete(usr *User) error {
	return us.DeleteByID(usr.ID)
}

// DeleteByID deletes one user from the database, by ID
func (us UserServiceStruct) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	result, err := us.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		//log.Println(err)
		return err
	}
	deleted, err := result.RowsAffected()
	if err != nil {
		//log.Println(err)
		return errors.New("failed to delete user")
	} else if deleted == 0 {
		//fmt.Println("deleted user ", deleted)
		return errors.New("failed to delete user")
	}

	return nil
}

// Delete All rows from users table
func (us UserServiceStruct) ResetUsers() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users`

	_, err := us.DB.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (us UserServiceStruct) Insert(user User) (User, error) {
	var usr User
	//ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	//defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	//log.Println("register hash password:", hashedPassword)
	if err != nil {
		return usr, err
	}
	//log.Println("Roll id:", user.RollId)
	if user.RollId == 0 {
		user.RollId = 3
	}
	//var newID int

	stmt := `INSERT INTO users (email, first_name, last_name, hashad_password, active, roll_id) values ($1, $2, $3, $4, $5, $6) 
	returning id, email, first_name, last_name, active, roll_id, updated_at, created_at`
	//log.Println("stmt:", stmt)
	err = us.DB.QueryRow(stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		0,
		user.RollId,
	).Scan(&usr.ID, &usr.Email, &usr.FirstName, &usr.LastName, &usr.Active, &usr.RollId, &usr.UpdatedAt, &usr.CreatedAt)

	if err != nil {
		//log.Println(err)
		return usr, err
	}

	return usr, nil
}

// ResetPassword is the method we will use to change a user's password.
func (us UserServiceStruct) ResetPassword(password string, usr *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	//log.Println("Password is:", password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set hashad_password = $1, password_changed_at= $2 where email = $3`
	_, err = us.DB.ExecContext(ctx, stmt, hashedPassword, time.Now(), usr.Email)
	if err != nil {
		return err
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (us UserServiceStruct) PasswordMatches(plainText string, usr *User) (bool, error) {
	//log.Println(u.Password)
	//log.Println(plainText)
	err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(plainText))
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
func (us UserServiceStruct) ConvertToProtoUser(user User) proto.User {
	return proto.User{
		ID:                int64(user.ID),
		Email:             user.Email,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Password:          user.Password,
		Active:            int32(user.Active),
		RollId:            int64(user.RollId),
		Roll:              user.Roll,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		UpdatedAt:         timestamppb.New(user.UpdatedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
func (us UserServiceStruct) ConvertToModelsUser(user *proto.User) User {
	return User{
		ID:                int(user.ID),
		Email:             user.Email,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Password:          user.Password,
		Active:            int(user.Active),
		RollId:            int(user.RollId),
		Roll:              user.Roll,
		PasswordChangedAt: user.PasswordChangedAt.AsTime(),
		UpdatedAt:         user.UpdatedAt.AsTime(),
		CreatedAt:         user.CreatedAt.AsTime(),
	}
}
