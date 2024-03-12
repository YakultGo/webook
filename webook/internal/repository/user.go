package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"context"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)
var _ UserRepository = (*CachedUserRepository)(nil)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	UpdateById(ctx context.Context, u domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
}
type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(d dao.UserDAO, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   d,
		cache: c,
	}
}
func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
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
func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id: u.Id,
	}, nil
}
func (r *CachedUserRepository) UpdateById(ctx context.Context, u domain.User) error {
	return r.dao.UpdateById(ctx, r.domainToEntity(u))
}

func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}
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
		_ = r.cache.Set(ctx, u)
	}()
	return u, nil
}

func (r *CachedUserRepository) domainToEntity(user domain.User) dao.User {
	return dao.User{
		Id:          user.Id,
		Email:       user.Email,
		Password:    user.Password,
		Phone:       user.Phone,
		NickName:    user.NickName,
		Description: user.Description,
	}
}
