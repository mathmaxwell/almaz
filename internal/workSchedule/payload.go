package workschedule

import (
	"demo/purpleSchool/configs"
	"demo/purpleSchool/internal/auth"
	"demo/purpleSchool/pkg/db"
	"time"
)

type ScheduleForDay struct {
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

type WorkScheduleDeps struct {
	*configs.Config
	ScheduleRepository *ScheduleRepository
	AuthHandler        *auth.AuthHandler
}
type WorkScheduleHandler struct {
	*configs.Config
	ScheduleRepository *ScheduleRepository
	AuthHandler        *auth.AuthHandler
}
type ScheduleRepository struct {
	DataBase *db.Db
}
type createWorkScheduleRequest struct {
	Token      string `json:"token" validate:"required"`
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
	Token string `json:"token" validate:"required"`
	ScheduleForDay
}
type getWorkScheduleRequest struct {
	ScheduleForDay
}
type getWorkScheduleResponse struct {
	Token              string           `json:"token" validate:"required"`
	EndDaySchedule     int              `json:"endDay"`
	EndMonthSchedule   int              `json:"endMonth"`
	EndYearSchedule    int              `json:"endYear"`
	StartMonthSchedule int              `json:"startMonth"`
	StartYearSchedule  int              `json:"startYear"`
	StartDaySchedule   int              `json:"startDay"`
	WorkSchedule       []ScheduleForDay `json:"workSchedule"`
}

type IWorkSchedule struct {
	StartDay     int              `json:"startDay"`
	StartMonth   int              `json:"startMonth"`
	StartYear    int              `json:"startYear"`
	EndDay       int              `json:"endDay"`
	EndMonth     int              `json:"endMonth"`
	EndYear      int              `json:"endYear"`
	WorkSchedule []ScheduleForDay `json:"workSchedule"`
}
type ITardinessHistory struct {
	Date         time.Time      `json:"date"`
	Id           string         `json:"id"`
	FullName     string         `json:"fullName"`
	Department   string         `json:"department"`
	Day          int            `json:"day"`
	Month        int            `json:"month"`
	Year         int            `json:"year"`
	EntryHour    int            `json:"entryHour"`
	EntryMinute  int            `json:"entryMinute"`
	EntryDay     int            `json:"entryDay"`
	EntryMonth   int            `json:"entryMonth"`
	EntryYear    int            `json:"entryYear"`
	ExitHour     int            `json:"exitHour"`
	ExitMinute   int            `json:"exitMinute"`
	ExitDay      int            `json:"exitDay"`
	ExitMonth    int            `json:"exitMonth"`
	ExitYear     int            `json:"exitYear"`
	WorkSchedule ScheduleForDay `json:"workSchedule"`
}
