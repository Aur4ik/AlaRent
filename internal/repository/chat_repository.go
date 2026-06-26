package repository

import (
	"time"

	"github.com/Aur4ik/AlaRent/internal/config"
	"github.com/Aur4ik/AlaRent/internal/models"
	"gorm.io/gorm"
)

func FindConversation(tenantID, landlordID, apartmentID uint) (*models.Conversation, error) {
	var conversation models.Conversation
	err := config.DB.
		Where("tenant_id = ? AND landlord_id = ? AND apartment_id = ?", tenantID, landlordID, apartmentID).
		First(&conversation).Error
	return &conversation, err
}

func CreateConversation(conversation *models.Conversation) error {
	return config.DB.Create(conversation).Error
}

func GetConversationByID(id uint) (*models.Conversation, error) {
	var conversation models.Conversation
	err := config.DB.
		Preload("Apartment").
		Preload("Apartment.Photos").
		Preload("Tenant").
		Preload("Landlord").
		First(&conversation, id).Error
	return &conversation, err
}

func GetUserConversations(userID uint) ([]models.Conversation, error) {
	var conversations []models.Conversation
	err := config.DB.
		Preload("Apartment").
		Preload("Apartment.Photos").
		Preload("Tenant").
		Preload("Landlord").
		Where("tenant_id = ? OR landlord_id = ?", userID, userID).
		Order("updated_at DESC").
		Find(&conversations).Error
	return conversations, err
}

func CreateMessage(message *models.Message) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(message).Error; err != nil {
			return err
		}

		return tx.Model(&models.Conversation{}).
			Where("id = ?", message.ConversationID).
			Update("updated_at", time.Now()).Error
	})
}

func GetMessageByID(id uint) (*models.Message, error) {
	var message models.Message
	err := config.DB.Preload("Sender").First(&message, id).Error
	return &message, err
}

func GetConversationMessages(conversationID uint) ([]models.Message, error) {
	var messages []models.Message
	err := config.DB.
		Preload("Sender").
		Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}
