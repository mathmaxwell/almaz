package main

import (
	"demo/purpleSchool/configs"
	"demo/purpleSchool/internal/auth"
	"demo/purpleSchool/internal/employees"
	"demo/purpleSchool/internal/link"
	"demo/purpleSchool/internal/messages"
	"demo/purpleSchool/internal/workSchedule"
	"demo/purpleSchool/pkg/cors"
	"demo/purpleSchool/pkg/db"
	"net/http"
)

func main() {

	conf := configs.LoadConfig()
	db := db.NewDB(conf)
	//repositories
	linkRepository := link.NewLinkRepository(db)
	router := http.NewServeMux()
	mux := cors.Cors(router)
	// Handler
	auth.NewAuthHandler(router, auth.AuthhandlerDeps{
		Config: conf,
	})
	employees.NewEmployeeHandler(router, employees.EmployeeshandlerDeps{
		Config: conf,
	})
	workschedule.NewWorkScheduleHandler(router, workschedule.WorkScheduleDeps{
		Config: conf,
	})
	messages.NewMessagesHandler(router, messages.MessagehandlerDeps{
		Config: conf,
	})

	link.NewLinkHandler(router, link.LinkhandlerDeps{
		LinkRepository: linkRepository,
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
