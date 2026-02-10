package announcements

import (
	"demo/almaz/configs"
	"demo/almaz/internal/auth"
	"demo/almaz/pkg/db"
)

type AnnouncementsRepository struct {
	DataBase *db.Db
}

type AnnouncementsRepositoryDeps struct {
	DataBase *db.Db
}

type Announcements struct {
	Id     string `json:"id"`
	Image  string `json:"image"`
	Uz     string `json:"uz"`
	Ru     string `json:"ru"`
	UzText string `json:"uzText"`
	RuText string `json:"ruText"`
}
type User struct {
	Login    string `gorm:"unique" json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Balance  int    `json:"balance"`
}

type AnnouncementsHandler struct {
	*configs.Config
	AnnouncementsRepository AnnouncementsRepository
	AuthHandler             *auth.AuthHandler
}

type AnnouncementshandlerDeps struct {
	*configs.Config
	AnnouncementsRepository *AnnouncementsRepository
	AuthHandler             *auth.AuthHandler
}

type DeleteAnnouncementsRequest struct {
	Token string `json:"token"`
	Id    string `json:"id"`
}
