package model

import (
	"github.com/sirupsen/logrus"

	proto "github.com/distribworks/dkron/v3/plugin/types"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUserFromProto(in *proto.UserModifyRequest, logger *logrus.Entry) *User {
	return &User{
		Username: in.Username,
		Password: in.Password,
	}
}

func (u *User) ToProto() *proto.UserModifyRequest {
	return &proto.UserModifyRequest{
		Username: u.Username,
	}
}
