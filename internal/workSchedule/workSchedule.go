package workschedule

import (
	"demo/purpleSchool/pkg/db"
	"demo/purpleSchool/pkg/req"
	"demo/purpleSchool/pkg/res"
	"net/http"
)

func NewScheduleRepository(dataBase *db.Db) *ScheduleRepository {
	return &ScheduleRepository{
		DataBase: dataBase,
	}
}

func NewWorkScheduleHandler(router *http.ServeMux, deps WorkScheduleDeps) {
	handler := &WorkScheduleHandler{
		Config: deps.Config,
	}
	router.HandleFunc("/workSchedule/updateWorkSchedule", handler.updateWorkSchedule())

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
