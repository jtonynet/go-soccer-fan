package database

import (
	"fmt"
	"log"

	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormConn struct {
	db *gorm.DB
}

func NewGormCom(cfg *config.Database) *GormConn {
	dsn := fmt.Sprintf(
		`host=%s user=%s password=%s dbname=%s port=%s sslmode=%s`,
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBname,
		cfg.Port,
		cfg.SSLmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	/*
	 TODO - Para fins de desenvolvimento. Quando ficar mais madura
	 removeremos os automigrates e adotaremos o golang migrate
	 https://github.com/golang-migrate/migrate
	*/
	if err := db.AutoMigrate(&model.Competition{}); err != nil {
		log.Fatalf("cannot automigrate competition: %v", err)
	}

	if err := db.AutoMigrate(&model.Team{}); err != nil {
		log.Fatalf("cannot automigrate team: %v", err)
	}

	if err := db.AutoMigrate(&model.Match{}); err != nil {
		log.Fatalf("cannot automigrate match: %v", err)
	}

	if err := db.AutoMigrate(&model.Fan{}); err != nil {
		log.Fatalf("cannot automigrate Fan: %v", err)
	}

	return &GormConn{db}
}

func (gc *GormConn) GetDB() *gorm.DB {
	return gc.db
}
