package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewConnection(db_host, db_port, db_name, db_user, db_pass string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db_user, db_pass, db_host, db_port, db_name)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
