package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "modernc.org/sqlite"

	"Task-Scheduler/config"
)

func CheckDB() (*sql.DB, error) {

	dbfile := config.GetDBFile()

	var install bool
	if _, err := os.Stat(dbfile); err == nil {
		fmt.Println("Database file exists")
	} else if errors.Is(err, os.ErrNotExist) {
		install = true
		file, err := os.Create(dbfile)
		if err != nil {
			fmt.Printf("Error creating database file: %s\n", err)
			return nil, err
		}
		err = file.Close()
		if err != nil {
			fmt.Printf("Error closing database file: %s\n", err)
			return nil, err
		}
		fmt.Println("Database created")
	}

	db, err := sql.Open("sqlite", dbfile)
	if err != nil {
		fmt.Printf("Error opening database file: %s\n", err)
		return nil, err
	}

	if install == true {
		_, err := db.Exec("CREATE TABLE scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT NOT NULL, title TEXT NOT NULL, comment TEXT, repeat TEXT);")
		if err != nil {
			fmt.Printf("Error creating table: %s\n", err)
			return nil, err
		}
		fmt.Println("Database table created")

		_, err = db.Exec("CREATE INDEX `date` ON `scheduler` (`date`);")
		if err != nil {
			fmt.Printf("Error creating index: %s\n", err)
			return nil, err
		}
		fmt.Println("Database table index created")
	}

	return db, nil
}
