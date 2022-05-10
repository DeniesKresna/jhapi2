package user

import (
	"context"

	"github.com/DeniesKresna/jhapi2/types"
)

type IRepository interface {
	GetUsers(ctx context.Context, search string) (student *types.Pagination, err error)
}
