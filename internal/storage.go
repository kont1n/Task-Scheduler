package internal

import (
	"database/sql"
	"strconv"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return Storage{db: db}
}

func (s *Storage) GetAllTasks() ([]Task, error) {
	rows, err := s.db.Query("SELECT id,date,title,comment,repeat FROM scheduler")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Task
	for rows.Next() {
		t := Task{}
		err := rows.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) SearchTasks(parameter string) ([]Task, error) {
	query := "%" + parameter + "%"
	rows, err := s.db.Query("SELECT id,date,title,comment,repeat FROM scheduler WHERE title like :query OR comment like :query", sql.Named("query", query))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Task
	for rows.Next() {
		t := Task{}
		err := rows.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) SearchData(data string) ([]Task, error) {
	rows, err := s.db.Query("SELECT id,date,title,comment,repeat FROM scheduler WHERE date like :query", sql.Named("query", data))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Task
	for rows.Next() {
		t := Task{}
		err := rows.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) CreateTask(t Task) (string, error) {
	res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

func (s *Storage) ReadTask(id int) (Task, error) {
	row := s.db.QueryRow("SELECT id,date,title,comment,repeat FROM scheduler WHERE id = :id", sql.Named("id", id))

	t := Task{}
	err := row.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return Task{}, err
	}
	return t, nil
}

func (s *Storage) UpdateTask(t Task) error {
	i, err := strconv.Atoi(t.Id)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("UPDATE scheduler SET date = :date, title= :title, comment=:comment, repeat=:repeat WHERE id = :id",
		sql.Named("id", i),
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteTask(id int) error {
	_, err := s.db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		return err
	}
	return nil
}
