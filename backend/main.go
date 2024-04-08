package main

import (
	"fmt"
	"log"
	"sfdb/pkg/storage"
	pg "sfdb/pkg/storage/postgres"
	"time"
)

func main() {
	var store storage.Interface
	store, err := pg.New("postgres://user:password@database:5432/sf")
	if err != nil {
		log.Fatal(err)
	}

	labelRest := pg.Label{Name: "rest"}
	labelRest.ID, err = store.NewLabel(labelRest)
	if err != nil {
		log.Fatal(err)
	}

	labelWork := pg.Label{Name: "work"}
	labelWork.ID, err = store.NewLabel(labelWork)
	if err != nil {
		log.Fatal(err)
	}

	taskSleep := pg.Task{
		Opened:     time.Now().Unix(),
		AuthorID:   1,
		AssignedID: 1,
		Title:      "To sleep",
		Content:    "Sleep 8 hours minimum",
	}
	taskSleep.ID, err = store.NewTask(taskSleep)
	if err != nil {
		log.Fatal(err)
	}

	err = store.TaskAddLabel(taskSleep.ID, labelRest.ID)
	if err != nil {
		log.Fatal(err)
	}

	taskCode := pg.Task{
		Opened:     time.Now().Unix(),
		AuthorID:   1,
		AssignedID: 1,
		Title:      "To code",
		Content:    "Code 8 hours maximum",
	}
	taskCode.ID, err = store.NewTask(taskCode)
	if err != nil {
		log.Fatal(err)
	}

	err = store.TaskAddLabel(taskCode.ID, labelWork.ID)
	if err != nil {
		log.Fatal(err)
	}

	res, err := store.Tasks(0, "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("All tasks: %v\n", res)

	res, err = store.Tasks(1, "rest")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Tasks of user 1 with \"rest\" label: %v\n", res)

	err = store.DeleteTask(taskSleep.ID)
	if err != nil {
		log.Fatal(err)
	}

	res, err = store.Tasks(1, "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Tasks of user 1 after \"To Sleep\" deleting: %v\n", res)

	taskCode.AuthorID = 2
	err = store.UpdateTask(taskCode)
	if err != nil {
		log.Fatal(err)
	}

	res, err = store.Tasks(2, "work")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Tasks of user 2 with \"work\" label after author changed: %v\n", res)
}
