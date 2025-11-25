package employees

import (
	"demo/purpleSchool/configs"
	"time"
)

type EmployeeshandlerDeps struct {
	*configs.Config
}
type EmployeesHandler struct {
	*configs.Config
}
type GetEmployeesByIdRequest struct {
	Token string `json:"token" validate:"required"`
	Id    string `json:"id" validate:"required"`
}

type GetEmployeesRequest struct {
	Token     string `json:"token" validate:"required"`
	Page      int    `json:"page"`
	Count     int    `json:"count"`
	SortField string `json:"sortField"`
	SortAsc   bool   `json:"sortAsc"`
}
type IEmployeesResponse struct {
	Id                         string `json:"id"`
	Gender                     string `json:"gender"`
	Passport_series_and_number string `json:"passport_series_and_number"`
	PINFL                      string `json:"PINFL"`
	Full_name                  string `json:"full_name"`
	Image                      string `json:"image"` //"uploads/photo.jpg"
	Department                 string `json:"department"`
	Position                   string `json:"position"`
	Terminated                 bool   `json:"terminated"`
	On_probation               bool   `json:"on_probation"`
	On_vacation                bool   `json:"on_vacation"`
	On_sick_leave              bool   `json:"on_sick_leave"`
	On_a_business_trip         bool   `json:"on_a_business_trip"`
	Absence                    bool   `json:"absence"`
	Date_of_birth              int    `json:"date_of_birth"`
	Birth_month                int    `json:"birth_month"`
	Year_of_birth              int    `json:"year_of_birth"`
	Place_of_birth             string `json:"place_of_birth"`
	Nationality                string `json:"nationality"`
	Email                      string `json:"Email"`
	Phone_number               string `json:"phone_number"`
	Work_schedule              string `json:"work_schedule"`
}
type GetLateEmployeesRequest struct {
	Token       string `json:"token" validate:"required"`
	End_day     int    `json:"end_day" validate:"required"`
	End_month   int    `json:"end_month" validate:"required"`
	End_year    int    `json:"end_year" validate:"required"`
	Start_day   int    `json:"start_day" validate:"required"`
	Start_month int    `json:"start_month" validate:"required"`
	Start_year  int    `json:"start_year" validate:"required"`
}
type getLateEmployeesByIdRequest struct {
	Token       string `json:"token" validate:"required"`
	Id          string `json:"id" validate:"required"`
	End_month   int    `json:"endMonth" validate:"required"`
	End_year    int    `json:"endYear" validate:"required"`
	Start_month int    `json:"startMonth" validate:"required"`
	Start_year  int    `json:"startYear" validate:"required"`
}
type GetLateEmployeesResponse struct {
	Terminated         int `json:"terminated"`
	On_probation       int `json:"on_probation"`
	Active_employees   int `json:"active_employees"`
	On_vacation        int `json:"on_vacation"`
	On_sick_leave      int `json:"on_sick_leave"`
	On_a_business_trip int `json:"on_a_business_trip"`
	Absence            int `json:"absence"`
	Total_employees    int `json:"total_employees"`
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
type IWorkSchedule struct {
	StartDay     int                   `json:"startDay"`
	StartMonth   int                   `json:"startMonth"`
	StartYear    int                   `json:"startYear"`
	EndDay       int                   `json:"endDay"`
	EndMonth     int                   `json:"endMonth"`
	EndYear      int                   `json:"endYear"`
	WorkSchedule []IWorkScheduleForDay `json:"workSchedule"`
}

type ITardinessHistory struct {
	Date         time.Time           `json:"date"`
	Id           string              `json:"id"`
	FullName     string              `json:"fullName"`
	Department   string              `json:"department"`
	Day          int                 `json:"day"`
	Month        int                 `json:"month"`
	Year         int                 `json:"year"`
	EntryHour    int                 `json:"entryHour"`
	EntryMinute  int                 `json:"entryMinute"`
	EntryDay     int                 `json:"entryDay"`
	EntryMonth   int                 `json:"entryMonth"`
	EntryYear    int                 `json:"entryYear"`
	ExitHour     int                 `json:"exitHour"`
	ExitMinute   int                 `json:"exitMinute"`
	ExitDay      int                 `json:"exitDay"`
	ExitMonth    int                 `json:"exitMonth"`
	ExitYear     int                 `json:"exitYear"`
	WorkSchedule IWorkScheduleForDay `json:"workSchedule"`
}
type IEmployeesCountRequest struct {
	Token string `json:"token" validate:"required"`
	Day   int    `json:"day" validate:"required"`
	Month int    `json:"month" validate:"required"`
	Year  int    `json:"year" validate:"required"`
}
type IEmployeesCountResponse struct {
	Terminated         int `json:"terminated"`
	On_probation       int `json:"on_probation"`
	Active_employees   int `json:"active_employees"`
	On_vacation        int `json:"on_vacation"`
	On_sick_leave      int `json:"on_sick_leave"`
	On_a_business_trip int `json:"on_a_business_trip"`
	Absence            int `json:"absence"`
	Total_employees    int `json:"total_employees"`
}
