package access

import (
	"github.com/mrVoldemar/crm_backend/services/auth/internal/service"
	desc "github.com/mrVoldemar/crm_backend/services/auth/pkg/access_v1"
)

type Implementation struct {
	desc.UnimplementedAccessV1Server
	accessService service.AccessService
}

func NewImplementation(accessService service.AccessService) *Implementation {
	return &Implementation{
		accessService: accessService,
	}
}
