package workschedule

import (
	"demo/purpleSchool/configs"
)

type WorkScheduleDeps struct {
	*configs.Config
}
type WorkScheduleHandler struct {
	*configs.Config
}
type GetWorkScheduleForMonthRequest struct {
	Token              string `json:"token" validate:"required"`
	Id                 string `json:"id" validate:"required"`
	StartMonthSchedule int    `json:"startMonthSchedule"`
	StartYearSchedule  int    `json:"startYearSchedule"`
}

type GetWorkScheduleForMonthResponse struct {
	StartDay     int                   `json:"startDay"`
	StartMonth   int                   `json:"startMonth"`
	StartYear    int                   `json:"startYear"`
	EndDay       int                   `json:"endDay"`
	EndMonth     int                   `json:"endMonth"`
	EndYear      int                   `json:"endYear"`
	WorkSchedule []IWorkScheduleForDay `json:"workSchedule"`
}
type IWorkScheduleForDay struct {
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

type updateWorkScheduleRequest struct {
	Token      string `json:"token" validate:"required"`
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
