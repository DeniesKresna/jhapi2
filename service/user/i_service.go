package user

import (
	"context"

	"github.com/DeniesKresna/jhapi2/types"
)

type IService interface {
	GetUsers(ctx context.Context, search string) (user *types.Pagination, err error)
}
