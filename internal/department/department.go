package department

import (
	"demo/purpleSchool/pkg/db"
	"demo/purpleSchool/pkg/req"
	"demo/purpleSchool/pkg/res"
	"demo/purpleSchool/pkg/token"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

func NewDeportamentRepository(dataBase *db.Db) *DeportamentRepository {
	return &DeportamentRepository{
		DataBase: dataBase,
	}
}
func NewDeportamentHandler(router *http.ServeMux, deps DepartmenthandlerDeps) *Departmenthandler {
	handler := &Departmenthandler{
		Config:                deps.Config,
		DeportamentRepository: deps.DeportamentRepository,
	}
	router.HandleFunc("/department/createDepartment", handler.createDepartment())
	router.HandleFunc("/department/getDepartment", handler.getDepartment())
	return handler
}
func (handler *Departmenthandler) createDepartment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[DepartmentRequest](&w, r)
		if err != nil {
			return
		}

		// Проверяем, есть ли уже департамент с таким именем
		var existing Department
		err = handler.DeportamentRepository.DataBase.Where("name = ?", body.Name).First(&existing).Error
		if err == nil {
			// запись найдена → ошибка
			res.Json(w, map[string]string{
				"error": "департамент с таким именем уже существует",
			}, 400)
			return
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			// ошибка базы данных
			res.Json(w, map[string]string{
				"error": "ошибка базы данных",
			}, 500)
			return
		}

		// Создаем новый департамент
		newDeportament := Department{
			Name: body.Name,
			Id:   token.CreateId(),
		}
		if err := handler.DeportamentRepository.DataBase.Create(&newDeportament).Error; err != nil {
			res.Json(w, map[string]string{
				"error": "не удалось создать департамент",
			}, 500)
			return
		}

		res.Json(w, newDeportament, 200)
	}
}
func (handler *Departmenthandler) getDepartment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := req.HandleBody[getDepartmentRequest](&w, r)
		if err != nil {
			return
		}
		var data []Department
		err = handler.DeportamentRepository.DataBase.Find(&data).Error
		if err != nil {
			res.Json(w, map[string]string{
				"error": "ошибка базы данных",
			}, 500)
			return
		}

		res.Json(w, data, 200)
	}
}
