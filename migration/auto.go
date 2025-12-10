package main

import (
	"demo/purpleSchool/internal/auth"
	"demo/purpleSchool/internal/employees"
	workschedule "demo/purpleSchool/internal/workSchedule"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&auth.LoginResponse{})
	db.AutoMigrate(&employees.Employee{})
	db.AutoMigrate(&employees.EmployeeStatus{})
	db.AutoMigrate(&workschedule.IWorkScheduleForDay{})
}
