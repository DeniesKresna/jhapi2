package user

import (
	"context"

	"github.com/DeniesKresna/jhapi2/repository/cache"
	"github.com/DeniesKresna/jhapi2/repository/sql"
	"github.com/DeniesKresna/jhapi2/types"
	"github.com/DeniesKresna/jhapi2/utils"
	"gorm.io/gorm"
)

const table = "users"

type RepositoryImpl struct {
	db        *gorm.DB
	cache     *cache.IRepository
	utilities utils.IUtils
}

func Provide(db *gorm.DB, cache *cache.IRepository, utilities utils.IUtils) IRepository {
	return &RepositoryImpl{
		db:        db,
		cache:     cache,
		utilities: utilities,
	}
}

func (r *RepositoryImpl) GetUsers(ctx context.Context, search string) (paginationData *types.Pagination, err error) {
	var users []types.User
	r.db.Scopes(sql.Paginate(users, paginationData, r.db)).Where("name like ?", "%"+search+"%").Find(&users)
	paginationData.Rows = users
	return
}
