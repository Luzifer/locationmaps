package main

import (
	"fmt"
	"strings"
)

var users *userDatabase

func loadUserDB() error {
	err := couchConn.Get("users", &users)
	if err != nil {
		if strings.Contains(err.Error(), "Code: 404") {
			users = &userDatabase{ID: "users", Users: make(map[string]user)}
			saveUserDB()
		} else {
			return fmt.Errorf("Unable to load userdb: %s\n", err)
		}
	}

	return nil
}

func saveUserDB() error {
	res, err := couchConn.Save(users)
	if err != nil {
		return err
	}

	if !res.Ok {
		return fmt.Errorf("Unable to store userdb: %s\n", res.Reason)
	}

	users.Rev = res.Rev
	return nil
}
