package access

import (
	"context"
	descAccess "github.com/mrVoldemar/crm_backend/services/auth/pkg/access_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Check(ctx context.Context, req *descAccess.CheckRequest) (*emptypb.Empty, error) {
	err := i.accessService.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil

}
