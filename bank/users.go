package bank

import "time"

func CreateUser(fullname string, username string, password string) SUser {
	id := NextUserID()
	return SUser{id, fullname, username, password, User, time.Now().UTC(), make([]SAccount, 0)}
}
