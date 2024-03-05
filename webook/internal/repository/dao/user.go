package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}
func (dao *UserDAO) Insert(ctx context.Context, u User) error {
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
			return ErrUserDuplicateEmail
		}
	}
	return err
}
func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) UpdateById(ctx context.Context, u User) error {
	u.UpdateTime = time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&u).Updates(User{
		NickName:    u.NickName,
		Birthday:    u.Birthday,
		Description: u.Description,
	}).Error
}
func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}

// User 直接对应数据库中的表结构
type User struct {
	Id          int64  `gorm:"primaryKey,autoIncrement"`
	Email       string `gorm:"unique"`
	NickName    string
	Birthday    time.Time
	Description string
	Password    string
	// 创建时间（毫秒）
	CreateTime int64
	// 更新时间
	UpdateTime int64
}
