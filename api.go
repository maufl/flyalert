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
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type NewUser struct {
	Name string
	Password string
}

func unmarshalBody(req *http.Request, v interface{}) error {
	bytes, err := ioutil.ReadAll(req.Body);
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, v); err != nil {
		return err
	}
	return nil
}

func register(req *http.Request) (int, string) {
	user := &NewUser{}
	err := unmarshalBody(req, user)
	if err != nil {
		return 500, err.Error()
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 500, err.Error()
	}
	database[user.Name] = &User{ Name: user.Name, Password: string(passwordHash) }
	return 200, "New user registered"
}

func authenticateUser(user, password string) bool {
	if user, ok := database[user]; ok {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
			return true
		}
	}
	return false
}

func getTriggers(user auth.User) (int, string) {
	if body, err := json.Marshal(database[string(user)].Triggers); err == nil {
		return 200, string(body)
	} else {
		return 500, err.Error()
	}
}

func addTrigger(user auth.User, req *http.Request) (int, string) {
	trigger := Trigger{}
	userStruct := database[string(user)]
	if err := unmarshalBody(req, &trigger); err != nil {
		return 500, err.Error()
	}
	userStruct.Triggers = append(userStruct.Triggers, trigger)
	return 200, "Trigger added"
}

func deleteTrigger(user auth.User, params martini.Params) (int, string) {
	index, err := strconv.ParseInt(params["id"], 0, 0)
	userStruct := database[string(user)]
	if err != nil {
		return 500, err.Error()
	}
	userStruct.Triggers = append(userStruct.Triggers[:index], userStruct.Triggers[index+1:]...)
	return 200, "Trigger removed"
}

func getEmails(user auth.User) (int, string) {
	if body, err := json.Marshal(database[string(user)].Emails); err == nil {
		return 200, string(body)
	} else {
		return 500, err.Error()
	}
}

func addEmail(user auth.User, req *http.Request) (int, string) {
	var email Email
	userStruct := database[string(user)]
	if err := unmarshalBody(req, &email); err != nil {
		return 500, err.Error()
	}
	userStruct.Emails = append(userStruct.Emails, email)
	return 200, "Email added"
}

func deleteEmail(user auth.User, params martini.Params) (int, string) {
	index, err := strconv.ParseInt(params["id"], 0, 0)
	userStruct := database[string(user)]
	if err != nil {
		return 500, err.Error()
	}
	userStruct.Emails = append(userStruct.Emails[:index], userStruct.Emails[index+1:]...)
	return 200, "Email removed"
}

func setupAPI() {
	m := martini.Classic()
	m.Group("/api/v1", func(r martini.Router) {
		r.Get("/triggers", getTriggers)
		r.Post("/triggers", addTrigger)
		r.Delete("/trigger/:id", deleteTrigger)
		r.Get("/emails", getEmails)
		r.Post("/emails", addEmail)
		r.Delete("/email/:id", deleteEmail)
	}, auth.BasicFunc(authenticateUser))
	m.Group("/api/v1", func(r martini.Router) {
		r.Post("/register", register)
	})
	m.Run()
}
