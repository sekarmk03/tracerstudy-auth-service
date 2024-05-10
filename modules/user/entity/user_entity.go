package entity

import (
	"time"
	"tracerstudy-auth-service/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const (
	UserTableName = "users"
)

type User struct {
	Id        uint64         `json:"id"`
	Name      string         `json:"name"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	RoleId    uint32         `json:"role_id"`
	CreatedAt time.Time      `gorm:"type:timestamptz;not_null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;not_null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func NewUser(id uint64, name, username, email, password string, roleId uint32) *User {
	return &User{
		Id:        id,
		Name:      name,
		Username:  username,
		Email:     email,
		Password:  password,
		RoleId:    roleId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *User) TableName() string {
	return UserTableName
}

func ConvertEntityToProto(u *User) *pb.User {
	return &pb.User{
		Id:        u.Id,
		Name:      u.Name,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		RoleId:    u.RoleId,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}
