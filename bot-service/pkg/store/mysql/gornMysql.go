package mysql

import (
	"fmt"
	"github.com/pkg/errors"
	config "github.com/serhiq/PhoneLinkerBot/internal/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Store struct {
	Db *gorm.DB
}

func New(s config.DBConfig) (*Store, error) {

	cfg := &gorm.Config{
		PrepareStmt: false,
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci&time_zone=UTC",
		s.Username,
		s.Password,
		s.Host,
		s.Port,
		s.DatabaseName)

	db, err := gorm.Open(mysql.Open(dsn), cfg)

	if err != nil {
		return nil, errors.Wrap(err, "unable initialize db session")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "unable get connection pull")
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "unable to ping DB")
	}

	return &Store{Db: db}, nil
}
