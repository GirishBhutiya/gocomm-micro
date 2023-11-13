package data

import (
	"testing"
	"time"
)

func TestUser_Insert(t *testing.T) {
	type fields struct {
		ID                int
		Email             string
		FirstName         string
		LastName          string
		Password          string
		Active            int
		RollId            int
		Roll              string
		PasswordChangedAt time.Time
		CreatedAt         time.Time
		UpdatedAt         time.Time
	}
	type args struct {
		user User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{name: "Insert user",
			fields:  fields{Email: "test@test.com", FirstName: "test F", LastName: "test L", Password: "secret", Active: 1, RollId: 1},
			args:    args{user: User{Email: "test@test.com", FirstName: "test F", LastName: "test L", Password: "secret", Active: 1, RollId: 1}},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:                tt.fields.ID,
				Email:             tt.fields.Email,
				FirstName:         tt.fields.FirstName,
				LastName:          tt.fields.LastName,
				Password:          tt.fields.Password,
				Active:            tt.fields.Active,
				RollId:            tt.fields.RollId,
				Roll:              tt.fields.Roll,
				PasswordChangedAt: tt.fields.PasswordChangedAt,
				CreatedAt:         tt.fields.CreatedAt,
				UpdatedAt:         tt.fields.UpdatedAt,
			}
			got, err := u.Insert(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("User.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}
