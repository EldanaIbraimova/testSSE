package model

import (
	"errors"
	"gorm.io/gorm"
)

type Message struct {
	Text   string `gorm:"size:100;" json:"text"`
	UserId uint32 `gorm:"size:100;" json:"uid"`
}

func (m *Message) SendMessage(db *gorm.DB) (*Message, error) {
	err := db.Debug().Create(&m).Error
	if err != nil {
		return &Message{}, err
	}
	return m, nil
}

func (m *Message) GetAllMessages(db *gorm.DB) ([]Message, error) {
	var err error
	var messages []Message
	err = db.Debug().Model(&Message{}).Find(&messages).Error
	if err != nil {
		return []Message{}, err
	}
	if len(messages) == 0 {
		return []Message{}, errors.New("departments not found")
	}
	return messages, nil
}
