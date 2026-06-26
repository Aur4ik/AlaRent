package service

import (
	"errors"
	"strings"

	"github.com/Aur4ik/AlaRent/internal/models"
	"github.com/Aur4ik/AlaRent/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrConversationForbidden = errors.New("you do not have access to this conversation")
	ErrConversationNotFound  = errors.New("conversation not found")
	ErrCannotMessageYourself = errors.New("you cannot start a conversation with yourself")
)

func StartConversation(userID, apartmentID uint) (*models.Conversation, error) {
	apartment, err := GetApartmentByID(apartmentID)
	if err != nil {
		return nil, err
	}

	if apartment.OwnerID == userID {
		return nil, ErrCannotMessageYourself
	}

	conversation, err := repository.FindConversation(userID, apartment.OwnerID, apartment.ID)
	if err == nil {
		return repository.GetConversationByID(conversation.ID)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	conversation = &models.Conversation{
		TenantID:    userID,
		LandlordID:  apartment.OwnerID,
		ApartmentID: apartment.ID,
	}
	if err := repository.CreateConversation(conversation); err != nil {
		return nil, err
	}

	return repository.GetConversationByID(conversation.ID)
}

func GetUserConversations(userID uint) ([]models.Conversation, error) {
	return repository.GetUserConversations(userID)
}

func GetConversationMessages(userID, conversationID uint) ([]models.Message, error) {
	if _, err := getConversationForUser(userID, conversationID); err != nil {
		return nil, err
	}

	return repository.GetConversationMessages(conversationID)
}

func SendMessage(userID, conversationID uint, text string) (*models.Message, error) {
	if _, err := getConversationForUser(userID, conversationID); err != nil {
		return nil, err
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return nil, errors.New("message text is required")
	}

	message := &models.Message{
		ConversationID: conversationID,
		SenderID:       userID,
		Text:           text,
	}
	if err := repository.CreateMessage(message); err != nil {
		return nil, err
	}

	return repository.GetMessageByID(message.ID)
}

func getConversationForUser(userID, conversationID uint) (*models.Conversation, error) {
	conversation, err := repository.GetConversationByID(conversationID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrConversationNotFound
	}
	if err != nil {
		return nil, err
	}

	if conversation.TenantID != userID && conversation.LandlordID != userID {
		return nil, ErrConversationForbidden
	}

	return conversation, nil
}
