package repository

import (
	"github.com/serhiq/PhoneLinkerBot/internal/repository/models"
)

type UserPhoneMappingRepository interface {
	InsertOrUpdate(mapping *models.UserPhoneMapping) error
	GetByChatID(chatID int64) (*models.UserPhoneMapping, error)
	GetChatIDByPhone(phoneNumber string) (*models.UserPhoneMapping, error)
	DeleteByChatID(chatID int64) error
}
