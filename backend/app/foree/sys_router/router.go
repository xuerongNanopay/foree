package sys_router

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	MysqlOpenConnection = `SHOW STATUS WHERE variable_name = ?`
)

func NewSystemRouter(db *sql.DB) *SystemRouter {
	return &SystemRouter{
		db: db,
	}
}

type SystemRouter struct {
	db *sql.DB
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to Foree"))
}

func (s *SystemRouter) mysqlConnection(w http.ResponseWriter, r *http.Request) {
	row := s.db.QueryRow(MysqlOpenConnection, "Threads_connected")

	var variableName string
	var value string

	err := row.Scan(&variableName, &value)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error: `%s`", err.Error())))
		return
	}

	w.Write([]byte(fmt.Sprintf("Mysql Threads_connected: `%s`", value)))
}

func (c *SystemRouter) RegisterRouter(router *mux.Router) {
	router.HandleFunc("/hello", hello).Methods("GET")
	router.HandleFunc("/mysql_connection", c.mysqlConnection).Methods("GET")
}
