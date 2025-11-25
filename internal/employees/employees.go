package employees

import (
	"demo/purpleSchool/pkg/req"
	"demo/purpleSchool/pkg/res"
	"demo/purpleSchool/pkg/token"
	"net/http"
	"time"
)

func NewEmployeeHandler(router *http.ServeMux, deps EmployeeshandlerDeps) {
	handler := &EmployeesHandler{
		Config: deps.Config,
	}
	router.HandleFunc("/employees/createEmployees", handler.createEmployees())
	router.HandleFunc("/employees/updateEmployees", handler.updateEmployees())
	router.HandleFunc("/employees/getEmployees", handler.getEmployees())
	router.HandleFunc("/employees/getLateEmployeesById", handler.getLateEmployeesById())
	router.HandleFunc("/employees/getEmployeesById", handler.getEmployeesById())
	router.HandleFunc("/employees/getLateEmployees", handler.getLateEmployees())
	router.HandleFunc("/employees/getEmployeesCount", handler.getEmployeesCount())
	
}
func (handler *EmployeesHandler) createEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := r.FormValue("token")
		if userToken == "" {
			res.Json(w, "unauthorized", 401)
			return
		}
		// gender := r.FormValue("gender")
		// Full_name := r.FormValue("full_name")
		// PINFL := r.FormValue("PINFL")
		// phone_number := r.FormValue("phone_number")
		// passport_series_and_number := r.FormValue("passport_series_and_number")
		// department := r.FormValue("department")
		// position := r.FormValue("position")
		// date_of_birth := r.FormValue("date_of_birth")
		// birth_month := r.FormValue("birth_month")
		// year_of_birth := r.FormValue("year_of_birth")
		// place_of_birth := r.FormValue("place_of_birth")
		// nationality := r.FormValue("nationality")
		// Email := r.FormValue("Email")
		file, header, err := r.FormFile("image")
		if err == nil && file != nil {
			res.Json(w, header, 200)
			return
		}
		res.Json(w, "error image", 401)
	}
}
func (handler *EmployeesHandler) updateEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := r.FormValue("token")
		employeeId := r.FormValue("id")
		if userToken == "" {
			res.Json(w, "unauthorized", 401)
			return
		}
		if employeeId == "" {
			res.Json(w, "unauthorized", 401)
			return
		}
		// gender := r.FormValue("gender")
		// Full_name := r.FormValue("full_name")
		// PINFL := r.FormValue("PINFL")
		// phone_number := r.FormValue("phone_number")
		// passport_series_and_number := r.FormValue("passport_series_and_number")
		// department := r.FormValue("department")
		// position := r.FormValue("position")
		// date_of_birth := r.FormValue("date_of_birth")
		// birth_month := r.FormValue("birth_month")
		// year_of_birth := r.FormValue("year_of_birth")
		// place_of_birth := r.FormValue("place_of_birth")
		// nationality := r.FormValue("nationality")
		// Email := r.FormValue("Email")
		file, header, err := r.FormFile("image")
		if err == nil && file != nil {
			res.Json(w, header, 200)
			return
		}
		res.Json(w, "error image", 401)
	}
}
func (handler *EmployeesHandler) getEmployeesById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetEmployeesByIdRequest](&w, r)
		if body.Token == "" {
			return
		} //проверка токена
		if err != nil {
			return
		}
		if body.Id == "" {
			return
		} //ищет по ид

		data := IEmployeesResponse{
			Id:                         token.CreateId(),
			Gender:                     "male",
			Passport_series_and_number: "AC1234567",
			PINFL:                      "12345678901234",
			Full_name:                  "Abdurahim Abdumalikov",
			Image:                      "image",
			Department:                 "Tad Industries",
			Position:                   "developer",
			Terminated:                 false,
			On_probation:               false,
			On_vacation:                false,
			On_sick_leave:              false,
			On_a_business_trip:         false,
			Absence:                    false,
			Date_of_birth:              9,
			Birth_month:                2,
			Year_of_birth:              2003,
			Place_of_birth:             "Tashkent",
			Nationality:                "uzbek",
			Email:                      "test123@gmail.com",
			Phone_number:               "+998(99)999-88-777",
			Work_schedule:              "полный день",
		}

		res.Json(w, data, 200)
	}
}
func (handler *EmployeesHandler) getEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetEmployeesRequest](&w, r)
		if body.Token == "" {
			return
		} //проверка токена
		if err != nil {
			return
		}
		//ищет за тех дней кто опоздал и сколько опоздал, график работы
		var data []IEmployeesResponse
		data1 := IEmployeesResponse{
			Id:                         "123",
			Gender:                     "male",
			Passport_series_and_number: "AC1234567",
			PINFL:                      "12345678901234",
			Full_name:                  "Abdurahim Abdumalikov",
			Image:                      "image",
			Department:                 "Tad Industries",
			Position:                   "developer",
			Terminated:                 false,
			On_probation:               false,
			On_vacation:                false,
			On_sick_leave:              false,
			On_a_business_trip:         false,
			Absence:                    false,
			Date_of_birth:              9,
			Birth_month:                2,
			Year_of_birth:              2003,
			Place_of_birth:             "Tashkent",
			Nationality:                "uzbek",
			Email:                      "test123@gmail.com",
			Phone_number:               "+998(99)999-88-777",
			Work_schedule:              "полный день",
		}
		data2 := IEmployeesResponse{
			Id:                         "321",
			Gender:                     "female",
			Passport_series_and_number: "AC7654321",
			PINFL:                      "09876543321232",
			Full_name:                  "test female",
			Image:                      "image",
			Department:                 "TADI",
			Position:                   "TADI developer",
			Terminated:                 false,
			On_probation:               false,
			On_vacation:                false,
			On_sick_leave:              false,
			On_a_business_trip:         false,
			Absence:                    false,
			Date_of_birth:              1,
			Birth_month:                3,
			Year_of_birth:              1998,
			Place_of_birth:             "Tashkent",
			Nationality:                "uzbek",
			Email:                      "test123@gmail.com",
			Phone_number:               "+998(99)999-88-777",
			Work_schedule:              "полный день",
		}

		data = append(data, data1)
		data = append(data, data2)
		res.Json(w, data, 200)
	}
}
func (handler *EmployeesHandler) getLateEmployeesById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[getLateEmployeesByIdRequest](&w, r)
		if body.Token == "" {
			return
		} //проверка токена
		if err != nil {
			return
		}
		//ищет за тех дней кто опоздал и сколько опоздал, график работы
		var data []ITardinessHistory
		data1 := ITardinessHistory{
			Date:        time.Now().AddDate(0, 0, 0),
			Id:          "Abdurahim Id",
			FullName:    "Abdurahim",
			Department:  "tad Industries",
			Day:         1,
			Month:       body.Start_month,
			Year:        body.Start_year,
			EntryHour:   9,
			EntryMinute: 20,
			EntryDay:    1,
			EntryMonth:  body.Start_month,
			EntryYear:   body.Start_year,
			ExitHour:    18,
			ExitMinute:  0,
			ExitDay:     1,
			ExitMonth:   body.Start_month,
			ExitYear:    body.Start_year,
			WorkSchedule: IWorkScheduleForDay{
				StartHour:  9,
				StartDay:   1,
				StartMonth: body.Start_month,
				StartYear:  body.Start_year,
				EndHour:    18,
				EndDay:     1,
				EndMonth:   body.Start_month,
				EndYear:    body.Start_year,
			},
		}
		data2 := ITardinessHistory{
			Date:        time.Now().AddDate(0, 0, 0),
			Id:          "Abdurahim Id",
			FullName:    "Abdurahim",
			Department:  "tad Industries",
			Day:         2,
			Month:       body.Start_month,
			Year:        body.Start_year,
			EntryHour:   8,
			EntryMinute: 55,
			EntryDay:    2,
			EntryMonth:  body.Start_month,
			EntryYear:   body.Start_year,
			ExitHour:    18,
			ExitMinute:  5,
			ExitDay:     2,
			ExitMonth:   body.Start_month,
			ExitYear:    body.Start_year,
			WorkSchedule: IWorkScheduleForDay{
				StartHour:  9,
				StartDay:   2,
				StartMonth: body.Start_month,
				StartYear:  body.Start_year,
				EndHour:    18,
				EndDay:     2,
				EndMonth:   body.Start_month,
				EndYear:    body.Start_year,
			},
		}
		data3 := ITardinessHistory{
			Date:        time.Now().AddDate(0, 0, 0),
			Id:          "Abdurahim Id",
			FullName:    "Abdurahim",
			Department:  "tad Industries",
			Day:         3,
			Month:       body.Start_month,
			Year:        body.Start_year,
			EntryHour:   9,
			EntryMinute: 5,
			EntryDay:    2,
			EntryMonth:  body.Start_month,
			EntryYear:   body.Start_year,
			ExitHour:    18,
			ExitMinute:  0,
			ExitDay:     3,
			ExitMonth:   body.Start_month,
			ExitYear:    body.Start_year,
			WorkSchedule: IWorkScheduleForDay{
				StartHour:  9,
				StartDay:   2,
				StartMonth: body.Start_month,
				StartYear:  body.Start_year,
				EndHour:    18,
				EndDay:     2,
				EndMonth:   body.Start_month,
				EndYear:    body.Start_year,
			},
		}
		data4 := ITardinessHistory{
			Date:        time.Now().AddDate(0, 0, 0),
			Id:          "Abdurahim Id",
			FullName:    "Abdurahim",
			Department:  "tad Industries",
			Day:         4,
			Month:       body.Start_month,
			Year:        body.Start_year,
			EntryHour:   10,
			EntryMinute: 0,
			EntryDay:    2,
			EntryMonth:  body.Start_month,
			EntryYear:   body.Start_year,
			ExitHour:    20,
			ExitMinute:  0,
			ExitDay:     4,
			ExitMonth:   body.Start_month,
			ExitYear:    body.Start_year,
			WorkSchedule: IWorkScheduleForDay{
				StartHour:  9,
				StartDay:   2,
				StartMonth: body.Start_month,
				StartYear:  body.Start_year,
				EndHour:    18,
				EndDay:     2,
				EndMonth:   body.Start_month,
				EndYear:    body.Start_year,
			},
		}
		data = append(data, data1)
		data = append(data, data2)
		data = append(data, data3)
		data = append(data, data4)
		res.Json(w, data, 200)
	}
}
func (handler *EmployeesHandler) getLateEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetLateEmployeesRequest](&w, r)
		if body.Token == "" {
			return
		} //проверка токена
		if err != nil {
			return
		}
		//ищет за тех дней кто опоздал и сколько опоздал, график работы
		var data []ITardinessHistory
		data1 := ITardinessHistory{
			Date:        time.Now().AddDate(0, 0, 0),
			Id:          "Abdurahim Id",
			FullName:    "Abdurahim",
			Department:  "tad Industries",
			Day:         body.Start_day,
			Month:       body.Start_month,
			Year:        body.Start_year,
			EntryHour:   9,
			EntryMinute: 20,
			EntryDay:    body.Start_day,
			EntryMonth:  body.Start_month,
			EntryYear:   body.Start_year,
			ExitHour:    18,
			ExitMinute:  0,
			ExitDay:     body.Start_day,
			ExitMonth:   body.Start_month,
			ExitYear:    body.Start_year,
			WorkSchedule: IWorkScheduleForDay{
				StartHour:  9,
				StartDay:   body.Start_day,
				StartMonth: body.Start_month,
				StartYear:  body.Start_year,
				EndHour:    18,
				EndDay:     body.Start_day,
				EndMonth:   body.Start_month,
				EndYear:    body.Start_year,
			},
		}
		data2 := ITardinessHistory{
			Date:        time.Now().AddDate(0, 0, 0),
			Id:          "employes Id",
			FullName:    "employes",
			Department:  "TADI",
			Day:         body.Start_day,
			Month:       body.Start_month,
			Year:        body.Start_year,
			EntryHour:   9,
			EntryMinute: 15,
			EntryDay:    body.Start_day,
			EntryMonth:  body.Start_month,
			EntryYear:   body.Start_year,
			ExitHour:    18,
			ExitMinute:  0,
			ExitDay:     body.Start_day,
			ExitMonth:   body.Start_month,
			ExitYear:    body.Start_year,
			WorkSchedule: IWorkScheduleForDay{
				StartHour:  9,
				StartDay:   body.Start_day,
				StartMonth: body.Start_month,
				StartYear:  body.Start_year,
				EndHour:    18,
				EndDay:     body.Start_day,
				EndMonth:   body.Start_month,
				EndYear:    body.Start_year,
			},
		}
		data = append(data, data1)
		data = append(data, data2)
		res.Json(w, data, 200)
	}
}
func (handler *EmployeesHandler) getEmployeesCount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[IEmployeesCountRequest](&w, r)
		if body.Token == "" {
			return
		} //проверка токена
		if err != nil {
			return
		}
		data := IEmployeesCountResponse{
			Terminated:         1,
			On_probation:       2,
			Active_employees:   3,
			On_vacation:        4,
			On_sick_leave:      5,
			On_a_business_trip: 6,
			Absence:            7,
			Total_employees:    100,
		}
		res.Json(w, data, 200)
	}
}
