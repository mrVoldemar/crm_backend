package register

import (
	"github.com/mrVoldemar/crm_backend/services/auth/internal/service"
	registerV1 "github.com/mrVoldemar/crm_backend/services/auth/pkg/register_v1"
)

type Implementation struct {
	registerV1.UnimplementedRegisterV1Server
	service service.RegisterService
}

func NewImplementation(service service.RegisterService) *Implementation {
	return &Implementation{
		service: service,
	}
}
