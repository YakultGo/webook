package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("账号/密码错误")
)

var _ UserService = (*UserServiceStruct)(nil)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, u domain.User) (domain.User, error)
	Edit(ctx context.Context, u domain.User) error
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}
type UserServiceStruct struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceStruct{
		repo: repo,
	}
}
func (svc *UserServiceStruct) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserServiceStruct) Login(ctx context.Context, u domain.User) (domain.User, error) {
	// 先查找用户
	dbU, err := svc.repo.FindByEmail(ctx, u.Email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 比较密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(dbU.Password), []byte(u.Password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return dbU, nil
}
func (svc *UserServiceStruct) Edit(ctx context.Context, u domain.User) error {
	return svc.repo.UpdateById(ctx, u)
}

func (svc *UserServiceStruct) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

func (svc *UserServiceStruct) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserNotFound) {
		return u, err
	}
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil {
		return domain.User{}, err
	}

	return svc.repo.FindByPhone(ctx, phone)
}
