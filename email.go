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
	"fmt"
	"net/smtp"
	"time"
)

type Email string

func createMessage(from, to string, notification Notification) []byte {
	message := `From: ` + from +`
To: ` + to + `
Subject: Flyalert!

It looks like you might be able to fly on ` + notification.Day.Format(time.RFC1123) + ` at ` + notification.Lat + ` ` + notification.Long + `.
http://www.openstreetmap.org/?mlat=` + notification.Lat + `&mlon=` + notification.Long
	return []byte(message)
}

func (e Email) SendNotification(n Notification) {
	fmt.Println("Sending notification ", n, " to address ", e)
	msg := createMessage(conf.Email.From, string(e), n)
	auth := smtp.PlainAuth("", conf.Email.Username, conf.Email.Password, conf.Email.Host)
	to := []string{ string(e) }
	err := smtp.SendMail(conf.Email.Server, auth, conf.Email.From, to, msg)
	if err != nil {
		fmt.Println("An error occured while sending an email: %v", err)
	}
}
