package user

import (
	"context"

	userSql "github.com/DeniesKresna/jhapi2/repository/sql/user"
	"github.com/DeniesKresna/jhapi2/types"
	"github.com/DeniesKresna/jhapi2/utils"
)

type ServiceImpl struct {
	userSql   userSql.IRepository
	utilities utils.IUtils
}

func Provide(userSql userSql.IRepository, utilities utils.IUtils) IService {
	return &ServiceImpl{
		userSql:   userSql,
		utilities: utilities,
	}
}

func (s *ServiceImpl) GetUsers(ctx context.Context, search string) (users *types.Pagination, err error) {
	if users, err = s.userSql.GetUsers(ctx, search); err != nil {
		return nil, err
	}
	return users, nil
}
