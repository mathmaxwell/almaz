package main

import (
	"demo/purpleSchool/configs"
	"demo/purpleSchool/internal/auth"
	"demo/purpleSchool/internal/department"
	"demo/purpleSchool/internal/employees"
	"demo/purpleSchool/internal/messages"
	"demo/purpleSchool/internal/workSchedule"
	"demo/purpleSchool/pkg/cors"
	"demo/purpleSchool/pkg/db"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDB(conf)
	router := http.NewServeMux()
	router.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	mux := cors.Cors(router)
	//register
	authRepo := auth.NewUserRepository(database)
	database.AutoMigrate(&auth.User{})
	authHandler := auth.NewAuthHandler(router, auth.AuthhandlerDeps{
		Config:         conf,
		AuthRepository: authRepo,
	})
	//employees
	employeeRepo := employees.NewEmployeeRepository(database)
	database.AutoMigrate(&employees.Employee{})
	database.AutoMigrate(&employees.EmployeeStatus{})
	employees.NewEmployeeHandler(router, employees.EmployeeshandlerDeps{
		Config:             conf,
		EmployeeRepository: employeeRepo,
		AuthHandler:        authHandler,
	})

	//work schedule
	scheduleRepo := workschedule.NewScheduleRepository(database)
	database.AutoMigrate(&workschedule.ScheduleForDay{})
	workschedule.NewWorkScheduleHandler(router, workschedule.WorkScheduleDeps{
		Config:             conf,
		ScheduleRepository: scheduleRepo,
		AuthHandler:        authHandler,
	})

	//message
	messages.NewMessagesHandler(router, messages.MessagehandlerDeps{
		Config: conf,
	})

	//department
	department.NewDeportamentHandler(router, department.DepartmenthandlerDeps{
		Config: conf,
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
