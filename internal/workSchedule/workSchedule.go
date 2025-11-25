package workschedule

import (
	"demo/purpleSchool/pkg/req"
	"demo/purpleSchool/pkg/res"
	"net/http"
)

func NewWorkScheduleHandler(router *http.ServeMux, deps WorkScheduleDeps) {
	handler := &WorkScheduleHandler{
		Config: deps.Config,
	}
	router.HandleFunc("/workSchedule/getWorkScheduleForMonth", handler.getWorkScheduleForMonth())
	router.HandleFunc("/workSchedule/updateWorkSchedule", handler.updateWorkSchedule())

}
func (handler *WorkScheduleHandler) getWorkScheduleForMonth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetWorkScheduleForMonthRequest](&w, r)
		if body.Token == "" {
			return
		} //проверка токена
		if err != nil {
			return
		}
		if body.Id == "" {
			return
		} //ищет по ид
		data1 := IWorkScheduleForDay{
			Id:         "bir",
			StartHour:  9,
			StartDay:   1,
			StartMonth: body.StartMonthSchedule,
			StartYear:  body.StartYearSchedule,
			EndHour:    18,
			EndDay:     1,
			EndMonth:   body.StartMonthSchedule,
			EndYear:    body.StartYearSchedule,
		}
		data2 := IWorkScheduleForDay{
			Id:         "ikki",
			StartHour:  9,
			StartDay:   2,
			StartMonth: body.StartMonthSchedule,
			StartYear:  body.StartYearSchedule,
			EndHour:    18,
			EndDay:     2,
			EndMonth:   body.StartMonthSchedule,
			EndYear:    body.StartYearSchedule,
		}
		data3 := IWorkScheduleForDay{
			Id:         "uch",
			StartHour:  9,
			StartDay:   3,
			StartMonth: body.StartMonthSchedule,
			StartYear:  body.StartYearSchedule,
			EndHour:    18,
			EndDay:     3,
			EndMonth:   body.StartMonthSchedule,
			EndYear:    body.StartYearSchedule,
		}
		data := GetWorkScheduleForMonthResponse{
			StartDay:   1,
			StartMonth: body.StartMonthSchedule,
			StartYear:  body.StartYearSchedule,
			EndDay:     30,
			EndMonth:   body.StartMonthSchedule,
			EndYear:    body.StartYearSchedule,
			WorkSchedule: []IWorkScheduleForDay{
				data1,
				data2,
				data3,
			},
		}

		res.Json(w, data, 200)
	}
}
func (handler *WorkScheduleHandler) updateWorkSchedule() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[updateWorkScheduleRequest](&w, r)

		if err != nil {
			return
		}

		//либо создает, либо меняет. если body.EndHour == body.EndHour==99 -> просто удалит данные. по дням работает, а не по ИД

		res.Json(w, body, 200)
	}
}
