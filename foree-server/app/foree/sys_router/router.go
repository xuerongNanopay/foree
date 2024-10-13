package sys_router

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	json_util "xue.io/go-pay/util/json"
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

func (s *SystemRouter) mysqlStats(w http.ResponseWriter, r *http.Request) {
	type MysqlStats struct {
		MaxOpenConnections int           `json:"maxOpenConnections"`
		OpenConnections    int           `json:"openConnections"`
		InUse              int           `json:"inUse"`
		Idle               int           `json:"idle"`
		WaitCount          int64         `json:"waitCount"`
		WaitDuration       time.Duration `json:"waitDuration"`
		MaxIdleClosed      int64         `json:"maxIdleClosed"`
		MaxIdleTimeClosed  int64         `json:"maxIdleTimeClosed"`
		MaxLifetimeClosed  int64         `json:"maxLifetimeClosed"`
	}

	dbStatus := s.db.Stats()

	mysqlStats := &MysqlStats{
		MaxOpenConnections: dbStatus.MaxOpenConnections,
		OpenConnections:    dbStatus.OpenConnections,
		InUse:              dbStatus.InUse,
		Idle:               dbStatus.Idle,
		WaitCount:          dbStatus.WaitCount,
		WaitDuration:       dbStatus.WaitDuration,
		MaxIdleClosed:      dbStatus.MaxIdleClosed,
		MaxIdleTimeClosed:  dbStatus.MaxIdleTimeClosed,
		MaxLifetimeClosed:  dbStatus.MaxLifetimeClosed,
	}

	json_util.SerializeToResponseWriter(w, http.StatusOK, mysqlStats)
}

func (c *SystemRouter) RegisterRouter(router *mux.Router) {
	router.HandleFunc("/hello", hello).Methods("GET")
	router.HandleFunc("/mysql_connection", c.mysqlConnection).Methods("GET")
	router.HandleFunc("/mysql_stats", c.mysqlStats).Methods("GET")

}
