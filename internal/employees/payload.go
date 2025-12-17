package employees

import (
	"demo/purpleSchool/configs"
	"demo/purpleSchool/internal/auth"
	"demo/purpleSchool/pkg/db"

	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	Id                         string `json:"id" gorm:"primaryKey"`
	Gender                     string `json:"gender"`
	Passport_series_and_number string `json:"passport_series_and_number"`
	PINFL                      string `json:"PINFL"`
	Full_name                  string `json:"full_name"`
	Image                      string `json:"image"`
	Department                 string `json:"department"`
	Position                   string `json:"position"`
	Date_of_birth              string `json:"date_of_birth"`
	Birth_month                string `json:"birth_month"`
	Year_of_birth              string `json:"year_of_birth"`
	Place_of_birth             string `json:"place_of_birth"`
	Nationality                string `json:"nationality"`
	Email                      string `json:"Email"`
	Phone_number               string `json:"phone_number"`
	Work_schedule              string `json:"work_schedule"`
	Accepted                   string `json:"accepted"`
}
type EmployeeStatus struct {
	gorm.Model
	Id         string `json:"id" validate:"required" gorm:"primaryKey"`
	EmployeeId string `json:"employeeId" validate:"required"`
	Status     string `json:"status"`
	StartDay   int    `json:"startDay"`
	StartMonth int    `json:"startMonth"`
	StartYear  int    `json:"startYear"`
	EndDay     int    `json:"endDay"`
	EndMonth   int    `json:"endMonth"`
	EndYear    int    `json:"endYear"`
}
type EmployeeRepository struct {
	DataBase *db.Db
}
type EmployeeshandlerDeps struct {
	*configs.Config
	EmployeeRepository *EmployeeRepository
	AuthHandler        *auth.AuthHandler
}

type EmployeesHandler struct {
	*configs.Config
	EmployeeRepository EmployeeRepository
	AuthHandler        *auth.AuthHandler
}
type GetEmployeesByIdRequest struct {
	Token string `json:"token" validate:"required"`
	Id    string `json:"id" validate:"required"`
}

type createStatusRequest struct {
	Token      string `json:"token" validate:"required"`
	EmployeeId string `json:"employeeId" validate:"required"`
	GetStatusByIdResponse
}
type GetStatusByIdResponse struct {
	Status     string `json:"status"`
	StartDay   int    `json:"startDay"`
	StartMonth int    `json:"startMonth"`
	StartYear  int    `json:"startYear"`
	EndDay     int    `json:"endDay"`
	EndMonth   int    `json:"endMonth"`
	EndYear    int    `json:"endYear"`
}

type GetEmployeesByStatusRequest struct {
	Token  string `json:"token" validate:"required"`
	Status string `json:"status" validate:"required"`
	Day    int    `json:"day" validate:"required"`
	Month  int    `json:"month" validate:"required"`
	Year   int    `json:"year" validate:"required"`
}
type GetEmployeesByStatusResponse struct {
	Ids []string `json:"ids"`
}

type GetEmployeesRequest struct {
	Token                      string `json:"token" validate:"required"`
	Page                       int    `json:"page"`
	Count                      int    `json:"count"`
	SortField                  string `json:"sortField"`
	SortAsc                    bool   `json:"sortAsc"`
	Gender                     string `json:"gender"`
	Passport_series_and_number string `json:"passport_series_and_number"`
	PINFL                      string `json:"PINFL"`
	Full_name                  string `json:"full_name"`
	Department                 string `json:"department"`
	Position                   string `json:"position"`
	Date_of_birth              string `json:"date_of_birth"`
	Birth_month                string `json:"birth_month"`
	Year_of_birth              string `json:"year_of_birth"`
	Place_of_birth             string `json:"place_of_birth"`
	Nationality                string `json:"nationality"`
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
