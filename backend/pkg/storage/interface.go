package storage

import pg "github.com/appstronomer/sf-hw-30.8.1/pkg/storage/postgres"

type Interface interface {
	NewTask(pg.Task) (int, error)
	Task(taskID int) (pg.Task, error)
	Tasks(authorID int, label string) ([]pg.Task, error)
	DeleteTask(taskID int) error
	UpdateTask(pg.Task) error
	NewLabel(pg.Label) (int, error)
	TaskAddLabel(taskID int, labelID int) error
}
