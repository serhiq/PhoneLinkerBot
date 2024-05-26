package userPhoneMapping

import (
	"github.com/serhiq/PhoneLinkerBot/internal/repository/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Repository struct {
	Db *gorm.DB
}

func New(Db *gorm.DB) *Repository {
	return &Repository{
		Db: Db,
	}
}

type GormDatabase struct {
	Db *gorm.DB
}

func CreateGorm(db *gorm.DB) *GormDatabase {
	return &GormDatabase{Db: db}
}

func (r *Repository) InsertOrUpdate(mapping *models.UserPhoneMapping) error {
	return r.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(mapping).Error
}

func (r *Repository) GetByChatID(chatID int64) (*models.UserPhoneMapping, error) {
	mapping := new(models.UserPhoneMapping)
	err := r.Db.Where("telegram_chat_id = ?", chatID).First(mapping).Error
	return mapping, err
}

func (r *Repository) GetChatIDByPhone(phone string) (*models.UserPhoneMapping, error) {
	mapping := new(models.UserPhoneMapping)
	err := r.Db.Where("phone_number = ?", phone).First(mapping).Error
	return mapping, err
}

func (r *Repository) DeleteByChatID(chatID int64) error {
	result := r.Db.Model(&models.UserPhoneMapping{}).Where("telegram_chat_id = ?", chatID).Update("deleted_at", time.Now())
	return result.Error
}
