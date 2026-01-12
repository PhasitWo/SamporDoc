package log

import (
	"SamporDoc/backend/infra/repository"
	"fmt"
	"time"
)

type Logger struct {
	repo          *repository.Repo
	correlationID *int64
}

func NewUnitOfLog(repo *repository.Repo) *Logger {
	id := time.Now().Unix()
	return &Logger{repo, &id}
}

func NewSingleLog(repo *repository.Repo) *Logger {
	return &Logger{repo, nil}
}

func (l *Logger) NewErrorAndLog(err error, action string) error {
	fmt.Println(action, repository.ERROR, l.correlationID, err.Error())
	l.repo.CreateLog(action, repository.ERROR, l.correlationID, err.Error())
	return err
}

func (l *Logger) Log(action string, data ...any) {
	fmt.Println(action, repository.SUCCESS, l.correlationID, data)
	l.repo.CreateLog(action, repository.SUCCESS, l.correlationID, data)
}
