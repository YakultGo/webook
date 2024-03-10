package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"context"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(d *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   d,
		cache: c,
	}
}
func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}
func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id: u.Id,
	}, nil
}
func (r *UserRepository) UpdateById(ctx context.Context, u domain.User) error {
	return r.dao.UpdateById(ctx, r.domainToEntity(u))
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}
	//if errors.Is(err, cache.ErrKeyNotExist) {
	//	// 缓存没有，从数据库中取
	//}
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = domain.User{
		Id:          ue.Id,
		Email:       ue.Email,
		NickName:    ue.NickName,
		Birthday:    ue.Birthday,
		Description: ue.Description,
	}
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 打日志，做监控
		}
	}()
	return u, nil
}

func (r *UserRepository) domainToEntity(user domain.User) dao.User {
	return dao.User{
		Id:          user.Id,
		Email:       user.Email,
		Password:    user.Password,
		Phone:       user.Phone,
		NickName:    user.NickName,
		Description: user.Description,
	}
}
