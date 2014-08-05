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
	"os"
	"encoding/json"
	"io/ioutil"
)
type User struct {
	Name string
	Password string
	Triggers []Trigger
	Notifications Notifications
	Emails []Email
}

type Database map[string]*User

func LoadDatabase(path string) (db Database) {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	e = json.Unmarshal(file, &db)
	if e != nil {
		fmt.Printf("Database error: %v\n", e)
		os.Exit(1)
	}
	return
}

func (db Database) Store(path string) {
	bytes, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		fmt.Printf("Could not serialize database! Changes are lost. Error: %v\n", err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(path, bytes, 0700)
	if err != nil {
		fmt.Printf("Could not write database! Changes are lost. Error: %v\n", err)
		os.Exit(1)
	}
	return
}
