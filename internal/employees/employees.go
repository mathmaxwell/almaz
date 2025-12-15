package employees

import (
	workschedule "demo/purpleSchool/internal/workSchedule"
	"demo/purpleSchool/pkg/db"
	"demo/purpleSchool/pkg/fields"
	"demo/purpleSchool/pkg/files"
	"demo/purpleSchool/pkg/req"
	"demo/purpleSchool/pkg/res"
	"demo/purpleSchool/pkg/token"
	"fmt"
	"net/http"
	"time"
)

func NewEmployeeRepository(dataBase *db.Db) *EmployeeRepository {
	return &EmployeeRepository{
		DataBase: dataBase,
	}
}

func NewEmployeeHandler(router *http.ServeMux, deps EmployeeshandlerDeps) *EmployeesHandler {
	handler := &EmployeesHandler{
		Config:             deps.Config,
		EmployeeRepository: *deps.EmployeeRepository,
		AuthHandler:        deps.AuthHandler,
	}
	router.HandleFunc("/employees/createEmployees", handler.createEmployees())
	router.HandleFunc("/employees/updateEmployees", handler.updateEmployees())
	router.HandleFunc("/employees/getEmployees", handler.getEmployees())
	router.HandleFunc("/employees/getLateEmployeesById", handler.getLateEmployeesById())
	router.HandleFunc("/employees/getEmployeesById", handler.getEmployeesById())
	router.HandleFunc("/employees/getLateEmployees", handler.getLateEmployees())
	router.HandleFunc("/employees/getEmployeesCount", handler.getEmployeesCount())
	router.HandleFunc("/employees/getStatusById", handler.getStatusById())
	router.HandleFunc("/employees/createStatus", handler.createStatus())
	router.HandleFunc("/employees/getEmployeesByStatus", handler.getEmployeesByStatus())
	return handler
}

func (handler *EmployeesHandler) createEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := r.FormValue("token")
		user, err := handler.AuthHandler.GetUserByToken(userToken)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		if user.UserRole != 1 {
			res.Json(w, "you are not admin", 403)
			return
		}
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			res.Json(w, "failed to parse form", http.StatusBadRequest)
			return
		}
		file, header, err := r.FormFile("image")
		if err != nil {
			res.Json(w, "image is not found", 400)
			return
		}
		photoPath, err := files.SaveFile(file, header)
		if err != nil {
			res.Json(w, "failed to save image", http.StatusInternalServerError)
			return
		}
		newEmployee := Employee{
			Id:                         token.CreateId(),
			Gender:                     r.FormValue("gender"),
			Full_name:                  r.FormValue("full_name"),
			PINFL:                      r.FormValue("PINFL"),
			Phone_number:               r.FormValue("phone_number"),
			Passport_series_and_number: r.FormValue("passport_series_and_number"),
			Department:                 r.FormValue("department"),
			Position:                   r.FormValue("position"),
			Date_of_birth:              r.FormValue("date_of_birth"),
			Birth_month:                r.FormValue("birth_month"),
			Year_of_birth:              r.FormValue("year_of_birth"),
			Place_of_birth:             r.FormValue("place_of_birth"),
			Nationality:                r.FormValue("nationality"),
			Email:                      r.FormValue("Email"),
			Image:                      photoPath,
		}
		if err := fields.ValidateFields(newEmployee, Employee{}); err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		if err := handler.EmployeeRepository.DataBase.Create(&newEmployee).Error; err != nil {
			res.Json(w, "db error", 500)
			return
		}
		response := map[string]interface{}{
			"message":     "employee created successfully",
			"employee_id": newEmployee.Id,
			"photo_url":   photoPath,
		}
		res.Json(w, response, 200)
	}
}
func (handler *EmployeesHandler) getEmployeesById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetEmployeesByIdRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		user, err := handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		if user.UserRole != 1 {
			res.Json(w, "you are not admin", 403)
			return
		}
		var employee Employee
		err = handler.EmployeeRepository.DataBase.Where("id = ?", body.Id).First(&employee).Error
		if err != nil {
			res.Json(w, "employee is not found", 400)
			return
		}
		res.Json(w, employee, 200)
	}
}
func (handler *EmployeesHandler) getEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetEmployeesRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		user, err := handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		if user.UserRole != 1 {
			res.Json(w, "you are not admin", 403)
			return
		}
		if body.Page <= 0 {
			body.Page = 1
		}
		if body.Count <= 0 {
			body.Count = 10
		}
		offset := (body.Page - 1) * body.Count
		sortOrder := "ASC"
		if !body.SortAsc {
			sortOrder = "DESC"
		}
		orderClause := ""
		if body.SortField == "date_of_birth" {
			orderClause = fmt.Sprintf(
				"make_date(year_of_birth::int, birth_month::int, date_of_birth::int) %s",
				sortOrder,
			)
		} else {
			sortField := "full_name"
			if body.SortField != "" {
				sortField = body.SortField
			}
			orderClause = fmt.Sprintf("%s %s", sortField, sortOrder)
		}
		var employees []Employee
		err = handler.EmployeeRepository.DataBase.
			Limit(body.Count).
			Offset(offset).
			Order(orderClause).
			Find(&employees).Error
		if err != nil {
			res.Json(w, "failed to get employees", 500)
			return
		}
		res.Json(w, employees, 200)
	}
}
func (handler *EmployeesHandler) getEmployeesCount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[IEmployeesCountRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		user, err := handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		if user.UserRole != 1 {
			res.Json(w, "you are not admin", 403)
			return
		}
		targetDate := body.Day + body.Month*100 + body.Year*1000
		db := handler.EmployeeRepository.DataBase
		var terminated, onProbation, onVacation, onSickLeave, onBusinessTrip, absence, total int64
		db.Model(&EmployeeStatus{}).
			Where("status = ?", "on_vacation").
			Where("start_day + start_month*100 + start_year*1000 <= ?", targetDate).
			Where("end_day + end_month*100 + end_year*1000 >= ?", targetDate).
			Count(&onVacation)

		db.Model(&EmployeeStatus{}).
			Where("status = ?", "on_sick_leave").
			Where("start_day + start_month*100 + start_year*1000 <= ?", targetDate).
			Where("end_day + end_month*100 + end_year*1000 >= ?", targetDate).
			Count(&onSickLeave)
		db.Model(&EmployeeStatus{}).
			Where("status = ?", "on_a_business_trip").
			Where("start_day + start_month*100 + start_year*1000 <= ?", targetDate).
			Where("end_day + end_month*100 + end_year*1000 >= ?", targetDate).
			Count(&onBusinessTrip)

		// db.Model(&Employee{}).Where("absence = ?", true).Count(&absence) //неявка
		db.Model(&Employee{}).Count(&total)
		data := IEmployeesCountResponse{
			Terminated:         int(terminated),
			On_probation:       int(onProbation),
			On_vacation:        int(onVacation),
			On_sick_leave:      int(onSickLeave),
			On_a_business_trip: int(onBusinessTrip),
			Absence:            int(absence),
			Total_employees:    int(total),
		}
		res.Json(w, data, 200)
	}
} //absence is not ready
func (handler *EmployeesHandler) updateEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := r.FormValue("token")
		user, err := handler.AuthHandler.GetUserByToken(userToken)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		if user.UserRole != 1 {
			res.Json(w, "you are not admin", 403)
			return
		}
		employeeId := r.FormValue("id")
		var employee Employee
		err = handler.EmployeeRepository.DataBase.Where("id = ?", employeeId).First(&employee).Error
		if err != nil {
			res.Json(w, "employee is not found", 400)
			return
		}
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			res.Json(w, "failed to parse form", http.StatusBadRequest)
			return
		}
		var photoPath string
		file, header, err := r.FormFile("image")
		if err == nil && file != nil {
			photoPath, err = files.SaveFile(file, header)
			if err != nil {
				res.Json(w, "failed to save image", 500)
				return
			}
		} else {
			photoPath = employee.Image
		}
		newEmployee := Employee{
			Id:                         employeeId,
			Gender:                     fields.GetOrDefault(r.FormValue("gender"), employee.Gender),
			Full_name:                  fields.GetOrDefault(r.FormValue("full_name"), employee.Full_name),
			PINFL:                      fields.GetOrDefault(r.FormValue("PINFL"), employee.PINFL),
			Phone_number:               fields.GetOrDefault(r.FormValue("phone_number"), employee.Phone_number),
			Passport_series_and_number: fields.GetOrDefault(r.FormValue("passport_series_and_number"), employee.Passport_series_and_number),
			Department:                 fields.GetOrDefault(r.FormValue("department"), employee.Department),
			Position:                   fields.GetOrDefault(r.FormValue("position"), employee.Position),
			Date_of_birth:              fields.GetOrDefault(r.FormValue("date_of_birth"), employee.Date_of_birth),
			Birth_month:                fields.GetOrDefault(r.FormValue("birth_month"), employee.Birth_month),
			Year_of_birth:              fields.GetOrDefault(r.FormValue("year_of_birth"), employee.Year_of_birth),
			Place_of_birth:             fields.GetOrDefault(r.FormValue("place_of_birth"), employee.Place_of_birth),
			Nationality:                fields.GetOrDefault(r.FormValue("nationality"), employee.Nationality),
			Email:                      fields.GetOrDefault(r.FormValue("Email"), employee.Email),
			Image:                      photoPath,
		}
		handler.EmployeeRepository.DataBase.Model(&employee).Updates(newEmployee)
		res.Json(w, "employee updated successfully", 200)
	}
}

func (handler *EmployeesHandler) createStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[createStatusRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		user, err := handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		if user.UserRole != 1 {
			res.Json(w, "you are not admin", 403)
			return
		}
		data := EmployeeStatus{
			Id:         token.CreateId(),
			EmployeeId: body.EmployeeId,
			Status:     body.Status,
			StartDay:   body.StartDay,
			StartMonth: body.StartMonth,
			StartYear:  body.StartYear,
			EndDay:     body.EndDay,
			EndMonth:   body.EndMonth,
			EndYear:    body.EndYear,
		}
		if err := handler.EmployeeRepository.DataBase.Create(&data).Error; err != nil {
			res.Json(w, "db error", 500)
			return
		}
		res.Json(w, body, 200)
	}
}
func (handler *EmployeesHandler) getStatusById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetEmployeesByIdRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 401)
			return
		}
		user, err := handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		if user.UserRole != 1 {
			res.Json(w, "you are not admin", 403)
			return
		}

		db := handler.EmployeeRepository.DataBase
		var employeeStatusHistory []EmployeeStatus
		db.Model(&EmployeeStatus{}).
			Where("employee_id = ?", body.Id).
			Find(&employeeStatusHistory)
		res.Json(w, employeeStatusHistory, 200)
	}
}

func (handler *EmployeesHandler) getEmployeesByStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetEmployeesByStatusRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 401)
			return
		}
		user, err := handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		if user.UserRole != 1 {
			res.Json(w, "you are not admin", 403)
			return
		}
		var ids []string
		targetDate := body.Day + body.Month*100 + body.Year*1000
		db := handler.EmployeeRepository.DataBase
		db.Model(&EmployeeStatus{}).Select("EmployeeId").
			Where("status = ?", body.Status).
			Where("start_day + start_month*100 + start_year*1000 <= ?", targetDate).
			Where("end_day + end_month*100 + end_year*1000 >= ?", targetDate).
			Find(&ids)
		res.Json(w, ids, 200)
	}
}

// finish
// work schedule !!!
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
		var data []workschedule.ITardinessHistory
		data1 := workschedule.ITardinessHistory{
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
			WorkSchedule: workschedule.ScheduleForDay{
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
		data2 := workschedule.ITardinessHistory{
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
			WorkSchedule: workschedule.ScheduleForDay{
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
		data3 := workschedule.ITardinessHistory{
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
			WorkSchedule: workschedule.ScheduleForDay{
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
		data4 := workschedule.ITardinessHistory{
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
			WorkSchedule: workschedule.ScheduleForDay{
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
		var data []workschedule.ITardinessHistory
		data1 := workschedule.ITardinessHistory{
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
			WorkSchedule: workschedule.ScheduleForDay{
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
		data2 := workschedule.ITardinessHistory{
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
			WorkSchedule: workschedule.ScheduleForDay{
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
