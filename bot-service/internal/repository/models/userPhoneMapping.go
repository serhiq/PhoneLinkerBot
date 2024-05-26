package models

import "time"

// UserPhoneMapping представляет соответствие между номерами телефонов и ID пользователей в Telegram.
type UserPhoneMapping struct {
	PhoneNumber string     `gorm:"column:phone_number;primaryKey;not null"`
	ChatID      int64      `gorm:"column:telegram_chat_id;not null"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

// TableName задает имя таблицы в базе данных.
func (UserPhoneMapping) TableName() string {
	return "user_phone_mapping"
}

func NewUserPhoneMapping(chatId int64, phone string) *UserPhoneMapping {
	p := new(UserPhoneMapping)
	p.ChatID = chatId
	p.PhoneNumber = phone
	return p
}
