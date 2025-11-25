package messages

import "demo/purpleSchool/configs"

type LoginResponse struct {
	Token    string `json:"token"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
type MessageHandler struct {
	*configs.Config
}

type MessagehandlerDeps struct {
	*configs.Config
}
type getUnreadMessagesRequest struct {
	Token string `json:"token" validate:"required"`
}
type getMessageWithUserRequest struct {
	Token string `json:"token" validate:"required"`
	Id    string `json:"id" validate:"required"`
}
type getUnreadMessagesResponse struct {
	Id           string `json:"id"`
	SenderId     string `json:"senderId"`
	RecipientIds string `json:"recipientIds"`
	Content      string `json:"content"`
	SentMinute   int    `json:"sentMinute"`
	SentHour     int    `json:"sentHour"`
	SentDay      int    `json:"sentDay"`
	SentMonth    int    `json:"sentMonth"`
	SentYear     int    `json:"sentYear"`
	IsRead       bool   `json:"isRead"`
}
type createResponse struct {
	Token        string `json:"token"`
	RecipientIds string `json:"recipientIds"`
	Content      string `json:"content"`
	SentMinute   int    `json:"sentMinute"`
	SentHour     int    `json:"sentHour"`
	SentDay      int    `json:"sentDay"`
	SentMonth    int    `json:"sentMonth"`
	SentYear     int    `json:"sentYear"`
}
