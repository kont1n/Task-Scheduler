package config

import (
	"os"
	"strconv"
)

const DBFile = "scheduler.db"
const Port = 7540

func GetDBFile() string {
	dbfile := DBFile
	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbfile = envFile
	}
	return dbfile
}

func GetServerPort() int {
	serverPort := Port
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) > 0 {
		if port, err := strconv.ParseInt(envPort, 10, 32); err == nil {
			if serverPort >= 0 {
				serverPort = int(port)
			}
		}
	}
	return serverPort
}
