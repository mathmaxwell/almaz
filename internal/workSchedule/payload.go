package workschedule

import (
	"demo/purpleSchool/configs"
	"demo/purpleSchool/internal/auth"
	"demo/purpleSchool/pkg/db"

	"gorm.io/gorm"
)

type WorkScheduleDeps struct {
	*configs.Config
	ScheduleRepository *ScheduleRepository
	AuthHandler        *auth.AuthHandler
}
type WorkScheduleHandler struct {
	*configs.Config
}
type ScheduleRepository struct {
	DataBase *db.Db
}

type IWorkScheduleForDay struct {
	gorm.Model
	Id         string `json:"id"`
	EmployeeId string `json:"employeeId"`
	StartHour  int    `json:"startHour"`
	StartDay   int    `json:"startDay"`
	StartMonth int    `json:"startMonth"`
	StartYear  int    `json:"startYear"`
	EndHour    int    `json:"endHour"`
	EndDay     int    `json:"endDay"`
	EndMonth   int    `json:"endMonth"`
	EndYear    int    `json:"endYear"`
}

type updateWorkScheduleRequest struct {
	Token      string `json:"token" validate:"required"`
	EmployeeId string `json:"employeeId"`
	Id         string `json:"id"`
	StartHour  int    `json:"startHour"`
	StartDay   int    `json:"startDay"`
	StartMonth int    `json:"startMonth"`
	StartYear  int    `json:"startYear"`
	EndHour    int    `json:"endHour"`
	EndDay     int    `json:"endDay"`
	EndMonth   int    `json:"endMonth"`
	EndYear    int    `json:"endYear"`
}
