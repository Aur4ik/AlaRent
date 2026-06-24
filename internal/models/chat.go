package models

import "gorm.io/gorm"

type Conversation struct {
	gorm.Model
	TenantID    uint      `gorm:"not null;index;uniqueIndex:idx_conversation_users_apartment"`
	LandlordID  uint      `gorm:"not null;index;uniqueIndex:idx_conversation_users_apartment"`
	ApartmentID uint      `gorm:"not null;index;uniqueIndex:idx_conversation_users_apartment"`
	Tenant      User      `gorm:"foreignKey:TenantID"`
	Landlord    User      `gorm:"foreignKey:LandlordID"`
	Apartment   Apartment `gorm:"foreignKey:ApartmentID"`
	Messages    []Message `gorm:"foreignKey:ConversationID"`
}

type Message struct {
	gorm.Model
	ConversationID uint         `gorm:"not null;index"`
	SenderID       uint         `gorm:"not null;index"`
	Text           string       `gorm:"type:text;not null"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID"`
	Sender         User         `gorm:"foreignKey:SenderID"`
}
