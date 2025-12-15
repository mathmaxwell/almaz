package workschedule

import (
	"demo/purpleSchool/pkg/db"
	"demo/purpleSchool/pkg/req"
	"demo/purpleSchool/pkg/res"
	"demo/purpleSchool/pkg/token"
	"net/http"
)

func NewScheduleRepository(dataBase *db.Db) *ScheduleRepository {
	return &ScheduleRepository{
		DataBase: dataBase,
	}
}

func NewWorkScheduleHandler(router *http.ServeMux, deps WorkScheduleDeps) {
	handler := &WorkScheduleHandler{
		Config:             deps.Config,
		ScheduleRepository: deps.ScheduleRepository,
		AuthHandler:        deps.AuthHandler,
	}
	router.HandleFunc("/workSchedule/getEmployeeWorkSchedule", handler.getEmployeeWorkSchedule())
	router.HandleFunc("/workSchedule/createWorkSchedule", handler.createWorkSchedule())
	router.HandleFunc("/workSchedule/updateWorkSchedule", handler.updateWorkSchedule())
}
func (handler *WorkScheduleHandler) createWorkSchedule() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[createWorkScheduleRequest](&w, r)
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
		newSchedule := ScheduleForDay{
			Id:         token.CreateId(),
			EmployeeId: body.EmployeeId,
			StartHour:  body.StartHour,
			StartDay:   body.StartDay,
			StartMonth: body.StartMonth,
			StartYear:  body.StartYear,
			EndHour:    body.EndHour,
			EndDay:     body.EndDay,
			EndMonth:   body.EndMonth,
			EndYear:    body.EndYear,
		}
		err = handler.ScheduleRepository.DataBase.Create(&newSchedule).Error
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		res.Json(w, newSchedule, 200)
	}
}

func (handler *WorkScheduleHandler) getEmployeeWorkSchedule() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[getWorkScheduleResponse](&w, r)
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
		schedule := []ScheduleForDay{}
		responce := getWorkScheduleResponse{
			EndDaySchedule:     body.EndDaySchedule,
			EndMonthSchedule:   body.EndMonthSchedule,
			EndYearSchedule:    body.EndYearSchedule,
			StartMonthSchedule: body.StartMonthSchedule,
			StartYearSchedule:  body.StartYearSchedule,
			StartDaySchedule:   body.StartDaySchedule,
			WorkSchedule:       schedule,
		}
		res.Json(w, responce, 200)
	}
}

func (handler *WorkScheduleHandler) updateWorkSchedule() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[updateWorkScheduleRequest](&w, r)
		if err != nil {
			return
		}
		//id это ид содрудника
		//либо создает, либо меняет. если body.EndHour == body.EndHour==99 -> просто удалит данные. по дням работает, а не по ИД

		res.Json(w, body, 200)
	}
}
