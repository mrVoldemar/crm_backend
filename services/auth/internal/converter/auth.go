package converter

import (
	"github.com/mrVoldemar/crm_backend/services/auth/internal/model"
	desc "github.com/mrVoldemar/crm_backend/services/auth/pkg/auth_v1"
)

func FromRequestToLoginModel(req *desc.LoginRequest) *model.UserLogin {
	return &model.UserLogin{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}
}
