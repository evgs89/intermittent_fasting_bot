package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DBConnector interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	Close()
	Connect() error
}

// Connector to mysql DB
type Connector struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DBname       string
	DBConnection *sql.DB
}

// NewConnection is a constructor for mysql Connector
func NewConnection(host string, port int, username, password, dbname string) *Connector {
	return &Connector{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		DBname:   dbname,
	}
}

// Connect to MySQL
func (c *Connector) Connect() error {
	dbconn, err := sql.Open("mysql", fmt.Sprintf("%v:%v@%v:%v/%v", c.Username, c.Password, c.Host, c.Port, c.DBname))
	c.DBConnection = dbconn
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// Close connection to MySQL
func (c *Connector) Close() {
	_ = c.DBConnection.Close()
}

// Exec query in MySQL
func (c *Connector) Exec(query string, args ...interface{}) (sql.Result, error) {
	if c.DBConnection == nil {
		_ = c.Connect()
	}
	return c.DBConnection.Exec(query, args)
}

// Query data from MySQL
func (c *Connector) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if c.DBConnection == nil {
		_ = c.Connect()
	}
	return c.DBConnection.Query(query, args)
}
