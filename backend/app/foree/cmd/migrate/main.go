package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	foree_config "xue.io/go-pay/app/foree/cmd/config"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	configPath := filepath.Join(ex, "../migrations")

	db, err := newMySQLStorage(mysqlCfg.Config{
		User:                 foree_config.Envs.DBUser,
		Passwd:               foree_config.Envs.DBPasswd,
		Addr:                 foree_config.Envs.DBAddr,
		DBName:               foree_config.Envs.DBName,
		MultiStatements:      true,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(configPath)
	fmt.Println(foree_config.Envs)

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("AAA")
	m, err := migrate.NewWithDatabaseInstance(
		configPath,
		foree_config.Envs.DBName,
		driver,
	)
	fmt.Println("AAA")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("AAA")
	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	} else if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("unknow cmd `%v`", cmd)
	}
}

func newMySQLStorage(cfg mysqlCfg.Config) (*sql.DB, error) {
	fmt.Println(cfg.FormatDSN())
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	return db, nil
}
