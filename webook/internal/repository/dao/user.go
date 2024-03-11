package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate = errors.New("邮箱冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)
var _ UserDAO = (*GORMUserDAO)(nil)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	UpdateById(ctx context.Context, u User) error
	FindById(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
}
type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}
func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	// 存毫秒
	now := time.Now().UnixMilli()
	u.UpdateTime = now
	u.CreateTime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConstraintErrNo = 1062
		if mysqlErr.Number == uniqueConstraintErrNo {
			// 邮箱冲突
			return ErrUserDuplicate
		}
	}
	return err
}
func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) UpdateById(ctx context.Context, u User) error {
	u.UpdateTime = time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&u).Updates(User{
		NickName:    u.NickName,
		Birthday:    u.Birthday,
		Description: u.Description,
		UpdateTime:  u.UpdateTime,
	}).Error
}
func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}
func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

// User 直接对应数据库中的表结构
type User struct {
	Id          int64     `gorm:"primaryKey,autoIncrement"`
	Email       string    `gorm:"unique;default:NULL"`
	Phone       string    `gorm:"unique;default:NULL"`
	NickName    string    `gorm:"default:NULL"`
	Birthday    time.Time `gorm:"default:NULL"`
	Description string    `gorm:"default:NULL"`
	Password    string    `gorm:"default:NULL"`
	// 单位毫秒
	CreateTime int64 `gorm:"default:NULL"`
	UpdateTime int64 `gorm:"default:NULL"`
}
