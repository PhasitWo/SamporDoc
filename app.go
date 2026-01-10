package main

import (
	"context"
	"fmt"

	"SamporDoc/backend/config"
	"SamporDoc/backend/model"
	"SamporDoc/backend/repository"
	"SamporDoc/backend/seed"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	ctx  context.Context
	repo *repository.Repo
}

func NewApp(ctx context.Context) *App {
	return &App{ctx: ctx}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// connect to the database
	dbFilePath, err := config.GetAppFilePath(config.DBFileName)
	if err != nil {
		fmt.Println(err)
		panic("Error getting database file path")
	}
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	// migrate
	db.AutoMigrate(&model.School{}, &model.Log{})

	// create repo instance
	a.repo = repository.NewRepo(ctx, db)

	// seeding
	schools, err := a.repo.GetAllSchools()
	if err != nil {
		fmt.Println("Error seeding database -> repo.getAllSchools")
		panic(err)
	}
	if len(schools) == 0 {
		fmt.Println("Start seeding...")
		for _, s := range seed.DefaultSchools {
			err = a.repo.CreateSchool(&s)
			if err != nil {
				fmt.Println("Error seeding database -> repo.CreateSchool")
				panic(err)
			}
		}
		fmt.Println("Finish seeding")
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
