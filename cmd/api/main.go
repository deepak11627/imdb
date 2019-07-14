package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/deepak11627/imdb/db"
	"github.com/deepak11627/imdb/handler"
	"github.com/deepak11627/imdb/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	addr        string
	defaultAddr = ":8080"
)

func init() {
	// environment vars
	_ = godotenv.Load(".env")
	_ = godotenv.Load("/etc/environment")

	// Flag values to override defaults
	flag.StringVar(&addr, "addr", defaultAddr, "Address at which to serve up api endpoints")
}
func main() {
	flag.Parse()
	// logger
	logger := logger.NewLogger(logger.Config{Output: os.Stdout})

	fmt.Println(os.Getenv("__MYSQL_DB_DSN"))
	// Connect to Database
	mysqlDB, err := sql.Open("mysql", os.Getenv("__MYSQL_DB_DSN"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer mysqlDB.Close()
	// verify the connection with the database
	if err := mysqlDB.Ping(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Create new Service with required dependencies
	s := handler.NewServer(addr, handler.NewHandler(db.NewDB(mysqlDB)))
	defer s.Close()

	log.Fatal(s.Open())
}
