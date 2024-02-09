package models

import (
	"authentication/cmd/db"
	"authentication/internal/util"
	"fmt"
	"strings"
	"testing"
)

func setup() UserService {

	db := db.ConnectToTestDB()

	return &UserServiceStruct{
		DB: db,
	}

}

var user User
var us UserService

func init() {
	us = setup()
	us.ResetUsers()
	user = CreateRandomUser()
}

func TestUser_Insert(t *testing.T) {

	type args struct {
		user        User
		userService UserService
	}
	tests := []struct {
		name    string
		args    args
		want    User
		isErr   bool
		wantErr map[string]string
	}{
		{name: "Insert user",
			args:    args{user: user, userService: us},
			want:    user,
			isErr:   false,
			wantErr: nil,
		},
		{name: "duplicate user",
			args:    args{user: user, userService: us},
			want:    user,
			isErr:   true,
			wantErr: map[string]string{"uniqueConstraint": "23505"},
		},
		{name: "without required fields user",
			args: args{user: func(u User) User {
				u.Email = ""
				return u
			}(user), userService: us},
			want:    user,
			isErr:   true,
			wantErr: map[string]string{"uniqueConstraint": "23505"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.userService.Insert(tt.args.user)

			if err != nil && !tt.isErr {
				if !strings.Contains(err.Error(), tt.wantErr["uniqueConstraint"]) {
					t.Errorf("User.Insert() error = %v\n, wantErr %v\n", err, tt.wantErr)
					return
				}
			}
			if !tt.isErr && !CompareUser(&got, &tt.want) {
				t.Errorf("User.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestUser_GetByEmail(t *testing.T) {

	type args struct {
		user        User
		userService UserService
	}
	tests := []struct {
		name    string
		args    args
		want    User
		isErr   bool
		wantErr string
	}{
		{name: "get correct User",
			args:    args{user: user, userService: us},
			want:    user,
			isErr:   false,
			wantErr: "",
		},
		{name: "Wrong email",
			args:    args{user: User{Email: "abc@example.com"}, userService: us},
			want:    user,
			isErr:   true,
			wantErr: "no rows in result set",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.userService.GetByEmail(tt.args.user.Email)

			if err != nil && tt.isErr {
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("User.GetByEmail() error = %v\n, wantErr %v\n", err, tt.wantErr)
					return
				}
			}
			if !tt.isErr && !CompareUser(got, &tt.want) {
				t.Errorf("User.GetByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestUser_GetOne(t *testing.T) {

	user = CreateRandomUser()
	u, err := us.Insert(user)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		user        User
		userService UserService
	}
	tests := []struct {
		name    string
		args    args
		want    User
		isErr   bool
		wantErr string
	}{
		{name: "get correct User",
			args:    args{user: u, userService: us},
			want:    u,
			isErr:   false,
			wantErr: "",
		},
		{name: "Wrong id",
			args:    args{user: User{ID: 0}, userService: us},
			want:    u,
			isErr:   true,
			wantErr: "no rows in result set",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.args.userService.GetOne(tt.args.user.ID)

			if err != nil && tt.isErr {
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("User.GetOne() error = %v\n, wantErr %v\n", err, tt.wantErr)
					return
				}
			}
			if !tt.isErr && !CompareUser(got, &tt.want) {
				t.Errorf("User.GetOne() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestUser_Update(t *testing.T) {

	user = CreateRandomUser()
	u, err := us.Insert(user)
	if err != nil {
		t.Fatal(err)
	}
	u.Email = util.RandomEmail()
	u.FirstName = util.RandStringBytesMaskImprSrcUnsafe(7)
	type args struct {
		user        User
		userService UserService
	}
	tests := []struct {
		name    string
		args    args
		want    User
		isErr   bool
		wantErr string
	}{
		{name: "Update correct User",
			args:    args{user: u, userService: us},
			want:    u,
			isErr:   false,
			wantErr: "",
		},
		{name: "Wrong id",
			args:    args{user: User{ID: 0}, userService: us},
			want:    u,
			isErr:   true,
			wantErr: "failed to update user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.args.userService.Update(&tt.args.user)

			if err != nil && tt.isErr {
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("User.Update() error = %v\n, wantErr %v\n", err, tt.wantErr)
					return
				}
			}

		})
	}
}
func TestUser_ResetPassword(t *testing.T) {

	user = CreateRandomUser()
	u, err := us.Insert(user)
	if err != nil {
		t.Fatal(err)
	}
	newPassword := fmt.Sprintf("%s%s", util.RandStringBytesMaskImprSrcUnsafe(7), "567!")

	type args struct {
		user        User
		userService UserService
	}
	tests := []struct {
		name    string
		args    args
		want    User
		isErr   bool
		wantErr string
	}{
		{name: "Update password",
			args:    args{user: u, userService: us},
			want:    u,
			isErr:   false,
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.args.userService.ResetPassword(newPassword, &tt.args.user)
			if err != nil {
				t.Fatal(err)
			}
			got, err := tt.args.userService.GetByEmail(tt.args.user.Email)
			if err != nil {
				t.Fatal(err)
			}
			ok, err := tt.args.userService.PasswordMatches(newPassword, got)
			if err != nil {
				t.Fatal(err)
			}
			if !ok {
				t.Errorf("User.ResetPassword() error = %v\n, wantErr %v\n", err, tt.wantErr)
			}

		})
	}
}
func TestUser_ValidateEmail(t *testing.T) {

	user = CreateRandomUser()
	u, err := us.Insert(user)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		user        User
		validate    bool
		userService UserService
	}
	tests := []struct {
		name    string
		args    args
		want    User
		isErr   bool
		wantErr string
	}{
		{name: "validate correct User",
			args:    args{user: u, userService: us, validate: true},
			want:    u,
			isErr:   false,
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.args.userService.ValidateEmail(&tt.args.user)

			if err != nil && tt.isErr {
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("User.GetOne() error = %v\n, wantErr %v\n", err, tt.wantErr)
					return
				}
			}

		})
	}
}
func TestUser_Delete(t *testing.T) {

	user = CreateRandomUser()
	u, err := us.Insert(user)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		user        User
		userService UserService
	}
	tests := []struct {
		name    string
		args    args
		isErr   bool
		wantErr string
	}{
		{name: "delete User",
			args:    args{user: u, userService: us},
			isErr:   false,
			wantErr: "",
		},
		{name: "Wrong id",
			args:    args{user: User{ID: 0}, userService: us},
			isErr:   true,
			wantErr: "failed to delete user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.args.userService.Delete(&tt.args.user)

			if err != nil && tt.isErr {
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("User.Delete() error = %v\n, wantErr %v\n", err, tt.wantErr)
					return
				}
			}

		})
	}
}
func CompareUser(got, want *User) bool {
	if got.Email != want.Email {
		return false
	}
	if got.FirstName != want.FirstName {
		return false
	}
	if got.LastName != want.LastName {
		return false
	}
	if got.RollId != want.RollId {
		return false
	}
	return true
}

// RandomEmail generates a random email
func CreateRandomUser() User {
	return User{
		Email:     util.RandomEmail(),
		FirstName: util.RandStringBytesMaskImprSrcUnsafe(6),
		LastName:  util.RandStringBytesMaskImprSrcUnsafe(6),
		Password:  fmt.Sprintf("%s%s", util.RandStringBytesMaskImprSrcUnsafe(7), "123!"),
		RollId:    1,
	}
}
