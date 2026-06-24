package repository

import (
	"github.com/Aur4ik/AlaRent/internal/config"
	"github.com/Aur4ik/AlaRent/internal/models"
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
		Preload("Tenant").
		Preload("Landlord").
		First(&conversation, id).Error
	return &conversation, err
}

func GetUserConversations(userID uint) ([]models.Conversation, error) {
	var conversations []models.Conversation
	err := config.DB.
		Preload("Apartment").
		Preload("Tenant").
		Preload("Landlord").
		Where("tenant_id = ? OR landlord_id = ?", userID, userID).
		Order("updated_at DESC").
		Find(&conversations).Error
	return conversations, err
}

func CreateMessage(message *models.Message) error {
	return config.DB.Create(message).Error
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
