package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User entity - sesuai dengan data di MongoDB
type User struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name      string               `bson:"name" json:"name"`
	Email     string               `bson:"email" json:"email"`
	Password  string               `bson:"password" json:"-"` // Never expose in JSON
	Role      UserRole             `bson:"role" json:"role"`
	Company   []primitive.ObjectID `bson:"company" json:"company"`
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time            `bson:"updatedAt" json:"updatedAt"`
	// Fields untuk forgot password (optional, bisa ditambah nanti)
	ResetPasswordToken   *string    `bson:"resetPasswordToken,omitempty" json:"-"`
	ResetPasswordExpires *time.Time `bson:"resetPasswordExpires,omitempty" json:"-"`
}

type UserRole string

const (
	RoleSuperAdmin UserRole = "SUPER_ADMIN"
	RoleAdmin      UserRole = "ADMIN"
	RoleClient     UserRole = "CLIENT"
)

func (r UserRole) IsValid() bool {
	switch r {
	case RoleSuperAdmin, RoleAdmin, RoleClient:
		return true
	}
	return false
}

// UserRepository interface
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, id primitive.ObjectID, user *User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	SetResetToken(ctx context.Context, email, token string, expires time.Time) error
	GetByResetToken(ctx context.Context, token string) (*User, error)
}
