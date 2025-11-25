package messages

import (
	"demo/purpleSchool/pkg/req"
	"demo/purpleSchool/pkg/res"
	"demo/purpleSchool/pkg/token"
	"net/http"
)

func NewMessagesHandler(router *http.ServeMux, deps MessagehandlerDeps) {
	handler := &MessageHandler{
		Config: deps.Config,
	}
	router.HandleFunc("/messages/getUnreadMessages", handler.getUnreadMessages())
	router.HandleFunc("/messages/getLastMessages", handler.getLastMessages())
	router.HandleFunc("/messages/getMessageWithUser", handler.getMessageWithUser())
	router.HandleFunc("/messages/readedMessage", handler.readedMessage())
	router.HandleFunc("/messages/createMessage", handler.createMessage())

}

func (handler *MessageHandler) createMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[createResponse](&w, r)
		//token == sendnerId
		if err != nil {
			return
		}
		if body.Token == "" {
			res.Json(w, "token", 400)
		}
		data := getUnreadMessagesResponse{
			Id:           token.CreateId(),
			SenderId:     body.Token,
			RecipientIds: body.RecipientIds,
			Content:      body.Content,
			SentMinute:   body.SentMinute,
			SentHour:     body.SentHour,
			SentDay:      body.SentDay,
			SentMonth:    body.SentMonth,
			SentYear:     body.SentYear,
			IsRead:       false,
		}
		res.Json(w, data, 200)
	}
}
func (handler *MessageHandler) readedMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[getMessageWithUserRequest](&w, r)
		//найдет по id, и меняет isRead false на true
		if err != nil {
			return
		}
		if body.Token == "" {
			res.Json(w, "token", 400)
		}

		res.Json(w, "update", 200)
	}
}
func (handler *MessageHandler) getUnreadMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[getUnreadMessagesRequest](&w, r)
		if err != nil {
			return
		}

		if body.Token == "" {
			res.Json(w, "token", 400)
		}
		//ищет token == recipientIds && isRead==false
		var data []getUnreadMessagesResponse
		msg1 := getUnreadMessagesResponse{
			Id:           "bir",
			SenderId:     "Sender-Id",
			RecipientIds: "RecipientIds===token",
			Content:      "msg1 Content",
			SentMinute:   0,
			SentHour:     16,
			SentDay:      18,
			SentMonth:    11,
			SentYear:     2025,
			IsRead:       false,
		}
		msg2 := getUnreadMessagesResponse{
			Id:           "bir",
			SenderId:     "Sender-Id",
			RecipientIds: "RecipientIds===token",
			Content:      "msg1 Content",
			SentMinute:   0,
			SentHour:     16,
			SentDay:      18,
			SentMonth:    11,
			SentYear:     2025,
			IsRead:       false,
		}
		data = append(data, msg1)
		data = append(data, msg2)

		res.Json(w, data, 200)
	}
}
func (handler *MessageHandler) getLastMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[getUnreadMessagesRequest](&w, r)
		if err != nil {
			return
		}

		if body.Token == "" {
			res.Json(w, "token", 400)
		}
		//ищет token == recipientIds && isRead==false
		var data []getUnreadMessagesResponse
		msg1 := getUnreadMessagesResponse{
			Id:           "bir",
			SenderId:     "tokken",
			RecipientIds: "74d144102322492e4efb7e8226dc6570",
			Content:      "получатель 74d144102322492e4efb7e8226dc6570",
			SentMinute:   0,
			SentHour:     16,
			SentDay:      19,
			SentMonth:    1,
			SentYear:     2025,
			IsRead:       true,
		}
		msg2 := getUnreadMessagesResponse{
			Id:           "ikki",
			SenderId:     "74d144102322492e4efb7e8226dc6570",
			RecipientIds: "token random",
			Content:      "отправитель 74d144102322492e4efb7e8226dc6570",
			SentMinute:   0,
			SentHour:     16,
			SentDay:      10,
			SentMonth:    12,
			SentYear:     2025,
			IsRead:       false,
		}
		msg3 := getUnreadMessagesResponse{
			Id:           "uch",
			SenderId:     "74d144102322492e4efb7e8226dc6570",
			RecipientIds: "random token",
			Content:      "отправитель 74d144102322492e4efb7e8226dc6570",
			SentMinute:   0,
			SentHour:     16,
			SentDay:      18,
			SentMonth:    11,
			SentYear:     2025,
			IsRead:       false,
		}
		msg4 := getUnreadMessagesResponse{
			Id:           "tort",
			SenderId:     "skdansfs",
			RecipientIds: "74d144102322492e4efb7e8226dc6570",
			Content:      "получатель 74d144102322492e4efb7e8226dc6570",
			SentMinute:   0,
			SentHour:     16,
			SentDay:      18,
			SentMonth:    11,
			SentYear:     2025,
			IsRead:       false,
		}
		msg5 := getUnreadMessagesResponse{
			Id:           "besh",
			SenderId:     "74d144102322492e4efb7e8226dc6570",
			RecipientIds: "random token",
			Content:      "отправитель 74d144102322492e4efb7e8226dc6570",
			SentMinute:   0,
			SentHour:     16,
			SentDay:      18,
			SentMonth:    11,
			SentYear:     2025,
			IsRead:       true,
		}
		msg6 := getUnreadMessagesResponse{
			Id:           "olti",
			SenderId:     "74d144102322492e4efb7e8226dc6570",
			RecipientIds: "random token",
			Content:      "отправитель 74d144102322492e4efb7e8226dc6570",
			SentMinute:   0,
			SentHour:     16,
			SentDay:      18,
			SentMonth:    11,
			SentYear:     2025,
			IsRead:       true,
		}
		data = append(data, msg1)
		data = append(data, msg2)
		data = append(data, msg3)
		data = append(data, msg4)
		data = append(data, msg5)
		data = append(data, msg6)
		data = append(data, msg6)
		data = append(data, msg6)
		data = append(data, msg6)
		data = append(data, msg6)
		data = append(data, msg6)
		data = append(data, msg6)
		data = append(data, msg6)
		data = append(data, msg6)
		data = append(data, msg6)
		data = append(data, msg6)
		data = append(data, msg6)
		data = append(data, msg6)
		data = append(data, msg6)
		data = append(data, msg6)
		data = append(data, msg6)
		res.Json(w, data, 200)
	}
}
func (handler *MessageHandler) getMessageWithUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//сначало самые новые сообщение
		body, err := req.HandleBody[getMessageWithUserRequest](&w, r)
		if err != nil {
			return
		}

		if body.Token == "" {
			res.Json(w, "token", 400)
		}
		if body.Id == "" {
			res.Json(w, "token", 400)
		}

		var data []getUnreadMessagesResponse
		msg1 := getUnreadMessagesResponse{
			Id:           token.CreateId(),
			SenderId:     body.Id,
			RecipientIds: "74d144102322492e4efb7e8226dc6570",
			Content:      token.CreateId(),
			SentMinute:   12,
			SentHour:     12,
			SentDay:      12,
			SentMonth:    12,
			SentYear:     2025,
			IsRead:       true,
		}
		msg2 := getUnreadMessagesResponse{
			Id:           token.CreateId(),
			SenderId:     "74d144102322492e4efb7e8226dc6570",
			RecipientIds: body.Id,
			Content:      token.CreateId(),
			SentMinute:   11,
			SentHour:     11,
			SentDay:      11,
			SentMonth:    11,
			SentYear:     2025,
			IsRead:       true,
		}
		msg3 := getUnreadMessagesResponse{
			Id:           token.CreateId(),
			SenderId:     body.Id,
			RecipientIds: "74d144102322492e4efb7e8226dc6570",
			Content:      token.CreateId(),
			SentMinute:   10,
			SentHour:     10,
			SentDay:      10,
			SentMonth:    10,
			SentYear:     2025,
			IsRead:       false,
		}
		msg4 := getUnreadMessagesResponse{
			Id:           token.CreateId(),
			SenderId:     "74d144102322492e4efb7e8226dc6570",
			RecipientIds: body.Id,
			Content:      token.CreateId(),
			SentMinute:   9,
			SentHour:     9,
			SentDay:      9,
			SentMonth:    9,
			SentYear:     2025,
			IsRead:       false,
		}

		data = append(data, msg1)
		data = append(data, msg2)
		data = append(data, msg3)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)
		data = append(data, msg4)

		res.Json(w, data, 200)
	}
}
