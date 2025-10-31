package register

import (
	"context"

	"github.com/mrVoldemar/crm_backend/services/auth/internal/model"
	registerV1 "github.com/mrVoldemar/crm_backend/services/auth/pkg/register_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Register(
	ctx context.Context,
	req *registerV1.RegisterRequest,
) (*emptypb.Empty, error) {
	user := &model.UserRegister{
		Email:      req.Email,
		Password:   req.Password,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Phone:      req.Phone,
		Position:   req.Position,
		Department: req.Department,
		AvatarURL:  req.AvatarUrl,
	}

	if req.HireDate != nil {
		user.HireDate = req.HireDate.AsTime()
	}

	if req.BirthDate != nil {
		user.BirthDate = req.BirthDate.AsTime()
	}

	_, err := i.service.Register(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &emptypb.Empty{}, nil
}
