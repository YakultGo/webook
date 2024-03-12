package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	smysql "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(t *testing.T) *sql.DB
		ctx     context.Context
		user    User
		wantErr error
		wantId  int64
	}{
		{
			name:    "插入成功",
			user:    User{},
			wantErr: nil,
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnResult(sqlmock.NewResult(3, 1))
				require.NoError(t, err)
				return mockDB
			},
		},
		{
			name:    "邮箱冲突",
			user:    User{},
			wantErr: ErrUserDuplicate,
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(&smysql.MySQLError{Number: 1062})
				require.NoError(t, err)
				return mockDB
			},
		},
		{
			name:    "数据库错误",
			user:    User{},
			wantErr: errors.New("数据库错误"),
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(errors.New("数据库错误"))
				require.NoError(t, err)
				return mockDB
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      tc.mock(t),
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				SkipDefaultTransaction: true,
				DisableAutomaticPing:   true,
			})
			d := NewUserDAO(db)
			err = d.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
