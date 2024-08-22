package internal

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func ValidateId(id string) (numberId int, err error) {

	if len(id) == 0 {
		err := errors.New("не указан идентификатор")
		return 0, err
	}

	numberId, err = strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return numberId, nil
}

func ValidateTask(task *Task) (err error) {
	var taskDate time.Time

	if task.Title == "" {
		err = errors.New("не указан заголовок задачи")
		return err
	}

	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	if taskDate, err = time.Parse("20060102", task.Date); err != nil {
		err = errors.New("дата представлена в формате, отличном от YYYYMMDD")
		return err
	}

	if taskDate.Before(time.Now()) && taskDate.Format("20060102") != time.Now().Format("20060102") {
		if task.Repeat == "" {
			fmt.Println("task date before now")
			task.Date = time.Now().Format("20060102")
		} else {
			newDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				err = errors.New("правило повторения указано в неправильном формате")
				return err
			}
			task.Date = newDate
		}
	}
	return nil
}
