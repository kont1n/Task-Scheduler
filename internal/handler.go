package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Tasks   []Task `json:"tasks,omitempty"`
	Id      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NextDateHandle(w http.ResponseWriter, r *http.Request) {
	var err error
	var nowTime time.Time
	var newDate string

	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if nowTime, err = time.Parse("20060102", now); err != nil {
		JSONError(w, err, http.StatusBadRequest)
		return
	}

	if _, err = time.Parse("20060102", date); err != nil {
		JSONError(w, err, http.StatusBadRequest)
		return
	}

	if len(repeat) == 0 {
		err = errors.New("не указано правило повторения")
		JSONError(w, err, http.StatusBadRequest)
		return
	}

	newDate, err = NextDate(nowTime, date, repeat)
	if err != nil {
		JSONError(w, err, http.StatusBadRequest)
		return
	}
	_, err = w.Write([]byte(newDate))
	if err != nil {
		JSONError(w, err, http.StatusInternalServerError)
		return
	}
}

func GetTasksHandle(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var answer Response
		var err error
		var tasks []Task

		search := r.URL.Query().Get("search")
		fmt.Println("")
		fmt.Println("Get task request. Search:", search)

		if search != "" {
			if res1, err := time.Parse("02.01.2006", search); err == nil {
				res2 := res1.Format(time.DateOnly)
				str := strings.ReplaceAll(res2, "-", "")

				tasks, err = store.SearchData(str)
				if err != nil {
					JSONError(w, err, http.StatusInternalServerError)
					return
				}
			} else {
				tasks, err = store.SearchTasks(search)
				if err != nil {
					JSONError(w, err, http.StatusInternalServerError)
					return
				}
			}
		} else {
			tasks, err = store.GetAllTasks()
			if err != nil {
				JSONError(w, err, http.StatusInternalServerError)
				return
			}
		}

		if len(tasks) == 0 {
			var resp []byte
			w.WriteHeader(http.StatusNotFound)
			resp, err = json.Marshal(map[string][]Task{
				"tasks": make([]Task, 0),
			})
			if err != nil {
				JSONError(w, err, http.StatusInternalServerError)
				return
			}
			fmt.Println("response=", string(resp))
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(resp)
			if err != nil {
				JSONError(w, err, http.StatusInternalServerError)
				return
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		answer = Response{Tasks: tasks}
		AnswerPrepare(answer, w)
	}
}

func CreateTaskHandle(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var answer Response
		var err error
		var id string
		var task Task
		var buf bytes.Buffer

		_, err = buf.ReadFrom(r.Body)
		if err != nil {
			JSONError(w, err, http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			JSONError(w, err, http.StatusBadRequest)
			return
		}
		fmt.Println("")
		fmt.Println("Create task request:", task)

		err = ValidateTask(&task)
		if err != nil {
			JSONError(w, err, http.StatusBadRequest)
			return
		}

		id, err = store.CreateTask(task)
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		answer = Response{Id: id}
		AnswerPrepare(answer, w)
	}
}

func ReadTaskHandle(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var answer Response
		var err error
		var id string
		var numberId int
		var task Task

		id = r.URL.Query().Get("id")

		fmt.Println("")
		fmt.Println("Read task with id=:", id)

		numberId, err = ValidateId(id)
		if err != nil {
			answer = Response{Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			AnswerPrepare(answer, w)
			return
		}

		task, err = store.ReadTask(numberId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				w.WriteHeader(http.StatusNotFound)
				answer = Response{Error: "задача не найдена"}
				AnswerPrepare(answer, w)
				return
			} else {
				JSONError(w, err, http.StatusInternalServerError)
				return
			}
		}

		answer.Id = task.Id
		answer.Date = task.Date
		answer.Title = task.Title
		answer.Comment = task.Comment
		answer.Repeat = task.Repeat
		w.WriteHeader(http.StatusOK)
		fmt.Println("answer:", answer)
		AnswerPrepare(answer, w)
	}
}

func UpdateTaskHandle(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var answer Response
		var err error
		var task Task
		var numberId int
		var buf bytes.Buffer

		_, err = buf.ReadFrom(r.Body)
		if err != nil {
			JSONError(w, err, http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			JSONError(w, err, http.StatusBadRequest)
			return
		}
		fmt.Println("")
		fmt.Println("Edit task request:", task)

		numberId, err = ValidateId(task.Id)
		if err != nil {
			answer = Response{Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			AnswerPrepare(answer, w)
			return
		}

		_, err = store.ReadTask(numberId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				w.WriteHeader(http.StatusNotFound)
				answer = Response{Error: "задача не найдена"}
				AnswerPrepare(answer, w)
				return
			} else {
				JSONError(w, err, http.StatusInternalServerError)
				return
			}
		}

		err = ValidateTask(&task)
		if err != nil {
			answer = Response{Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			AnswerPrepare(answer, w)
			return
		}

		err = store.UpdateTask(task)
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}
		fmt.Println("task updated")
		w.WriteHeader(http.StatusCreated)
		answer = Response{}
		AnswerPrepare(answer, w)
	}
}

func DeleteTaskHandle(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var answer Response
		var err error
		var numberId int

		id := r.URL.Query().Get("id")

		fmt.Println("")
		fmt.Println("Delete task with id=", id)

		numberId, err = ValidateId(id)
		if err != nil {
			answer = Response{Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			AnswerPrepare(answer, w)
			return
		}

		_, err = store.ReadTask(numberId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				w.WriteHeader(http.StatusNotFound)
				answer = Response{Error: "задача не найдена"}
				AnswerPrepare(answer, w)
				return
			} else {
				JSONError(w, err, http.StatusInternalServerError)
				return
			}
		}

		err = store.DeleteTask(numberId)
		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}
		fmt.Println("task deleted")
		answer = Response{}
		w.WriteHeader(http.StatusOK)
		AnswerPrepare(answer, w)
	}
}

func DoneTaskHandle(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var answer Response
		var err error
		var numberId int
		var task Task
		var id string

		id = r.URL.Query().Get("id")

		fmt.Println("")
		fmt.Println("Done task with id=", id)

		numberId, err = ValidateId(id)
		if err != nil {
			answer = Response{Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			AnswerPrepare(answer, w)
			return
		}

		task, err = store.ReadTask(numberId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				w.WriteHeader(http.StatusNotFound)
				answer = Response{Error: "задача не найдена"}
				AnswerPrepare(answer, w)
				return
			} else {
				JSONError(w, err, http.StatusInternalServerError)
				return
			}
		}

		if task.Repeat == "" {
			err = store.DeleteTask(numberId)
			if err != nil {
				JSONError(w, err, http.StatusInternalServerError)
				return
			}
			fmt.Println("task deleted as done")
		} else {
			newDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				JSONError(w, err, http.StatusInternalServerError)
				return
			}

			task.Date = newDate
			err = store.UpdateTask(task)
			if err != nil {
				JSONError(w, err, http.StatusInternalServerError)
				return
			}
			fmt.Println("task transfer")
		}
		answer = Response{}
		w.WriteHeader(http.StatusOK)
		AnswerPrepare(answer, w)
	}
}

func AnswerPrepare(answer Response, w http.ResponseWriter) {
	var resp []byte
	var err error

	resp, err = json.Marshal(answer)
	if err != nil {
		JSONError(w, err, http.StatusInternalServerError)
		return
	}
	fmt.Println("response=", string(resp))
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		JSONError(w, err, http.StatusInternalServerError)
		return
	}
}

func JSONError(w http.ResponseWriter, err error, code int) {
	var resp []byte
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	answer := Response{Error: err.Error()}
	resp, err = json.Marshal(answer)
	if err != nil {
		fmt.Println("Error marshalling response:", err)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("Error writing response:", err)
		return
	}
}
