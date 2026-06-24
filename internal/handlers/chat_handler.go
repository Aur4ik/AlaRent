package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Aur4ik/AlaRent/internal/dto"
	"github.com/Aur4ik/AlaRent/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var chatUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartConversation(c *gin.Context) {
	apartmentID, ok := parseApartmentID(c)
	if !ok {
		return
	}

	conversation, err := service.StartConversation(uint(c.GetInt("user_id")), apartmentID)
	if err != nil {
		writeChatError(c, err)
		return
	}

	c.JSON(http.StatusOK, conversation)
}

func GetConversations(c *gin.Context) {
	conversations, err := service.GetUserConversations(uint(c.GetInt("user_id")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, conversations)
}

func GetMessages(c *gin.Context) {
	conversationID, ok := parseConversationID(c)
	if !ok {
		return
	}

	messages, err := service.GetConversationMessages(uint(c.GetInt("user_id")), conversationID)
	if err != nil {
		writeChatError(c, err)
		return
	}

	c.JSON(http.StatusOK, messages)
}

func SendMessage(c *gin.Context) {
	conversationID, ok := parseConversationID(c)
	if !ok {
		return
	}

	var input dto.SendMessageRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := service.SendMessage(uint(c.GetInt("user_id")), conversationID, input.Text)
	if err != nil {
		writeChatError(c, err)
		return
	}

	c.JSON(http.StatusCreated, message)
}

func ConversationWebSocket(c *gin.Context) {
	conversationID, ok := parseConversationID(c)
	if !ok {
		return
	}

	userID := uint(c.GetInt("user_id"))
	if _, err := service.GetConversationMessages(userID, conversationID); err != nil {
		writeChatError(c, err)
		return
	}

	conn, err := chatUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		var input dto.SendMessageRequest
		if err := conn.ReadJSON(&input); err != nil {
			return
		}

		message, err := service.SendMessage(userID, conversationID, input.Text)
		if err != nil {
			_ = conn.WriteJSON(gin.H{"error": err.Error()})
			continue
		}

		if err := conn.WriteJSON(message); err != nil {
			return
		}
	}
}

func parseConversationID(c *gin.Context) (uint, bool) {
	conversationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation id"})
		return 0, false
	}

	return uint(conversationID), true
}

func writeChatError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrApartmentNotFound), errors.Is(err, service.ErrConversationNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrConversationForbidden), errors.Is(err, service.ErrCannotMessageYourself):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
