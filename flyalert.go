// Copyright (C) 2014 Felix Maurer
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>

package main

import (
	"time"
	"os"
	"os/signal"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var database Database

type EmailConfiguration struct {
	Server string
	From string
	Username string
	Password string
	Host string
}

type Configuration struct {
	ApiKey string
	Database string
	Interval time.Duration
	Email EmailConfiguration
}

// This struct contains global state and is also used in other places
var conf Configuration

func main() {
	// Load configuration
	file, e := ioutil.ReadFile(os.Args[1])
	if e != nil {
		fmt.Printf("Error while reading configuration file: %v\n", e)
		os.Exit(1)
	}
	e = json.Unmarshal(file, &conf)
	if e != nil {
		fmt.Printf("Configuration file is invalid: %v\n", e)
		os.Exit(1)
	}

	ticker := time.NewTicker(conf.Interval * time.Minute)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	database = LoadDatabase(conf.Database)

	processTriggers()
	go setupAPI()

	for {
		select {
		case <- ticker.C:
			processTriggers()
			clearNotifications()
		case s:= <-signals:
			fmt.Println("Got signal: ", s)
			fmt.Println("Saving database and shutting down, bye!")
			database.Store(conf.Database)
			return
		}
	}
}

func processTriggers() {
	for _, user := range database {
		for _, trigger := range user.Triggers {
			newNotifications := trigger.Process()
			processNewNotifications(user, newNotifications)
		}
	}
}

func processNewNotifications(user *User, newNotifications Notifications) {
	for _, notification := range newNotifications {
		if ! user.Notifications.Contain(notification) {
			for _, email := range user.Emails {
				email.SendNotification(notification)
			}
			user.Notifications = user.Notifications.Add(notification)
		}
	}
}

func clearNotifications() {
	for _, user := range database {
		newNotifications := Notifications{}
		for _, notification := range user.Notifications {
			// If the notificiation is yet to expire we keep it
			if notification.Day.After(time.Now()) {
				newNotifications = append(newNotifications, notification)
			}
		}
		user.Notifications = newNotifications
	}
}
