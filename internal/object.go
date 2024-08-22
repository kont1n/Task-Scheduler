package internal

type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type DB interface {
	GetAllTasks() ([]Task, error)
	SearchTasks(parameter string) ([]Task, error)
	SearchData(data string) ([]Task, error)
	CreateTask(t Task) (string, error)
	ReadTask(id int) (Task, error)
	UpdateTask(t Task) error
	DeleteTask(id int) error
}
