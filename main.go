package main

import (
	"./mysql"
	"./userdata"
	"github.com/evgs89/go-simplesettings"
)

var users = []*userdata.UserData{}

var Settings = simplesettings.NewSettingsFromFile("settings.ini")

func main() {
	host := Settings.Get("mysql", "host").ParseString()
	port := Settings.Get("mysql", "port").ParseInt()
	username := Settings.Get("mysql", "username").ParseString()
	password := Settings.Get("mysql", "password").ParseString()
	dbname := Settings.Get("mysql", "dbname").ParseString()
	c := mysql.NewConnection(host, port, username, password, dbname)
	userdata.DBCONN = c
	return
}
