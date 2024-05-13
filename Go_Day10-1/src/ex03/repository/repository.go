package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
	pb "ex02/data/go"
)

type Repository struct {
	db *gorm.DB
}

type Data struct {
	SessionId string
	Frequency float64
	Timestamp time.Time
}

func (r *Repository) Connect() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	if r.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	if err = r.db.AutoMigrate(&Data{}); err != nil {
		log.Fatalf("could not migrate db: %v", err)
	}
}

func (r *Repository) POST(message *pb.Message) {
	r.db.Create(&Data{
		SessionId: message.SessionId,
		Frequency: message.Frequency,
		Timestamp: message.Timestamp.AsTime(),
	})
}
