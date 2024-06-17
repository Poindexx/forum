package forum

import "database/sql"

var db *sql.DB
var sessionsMap = make(map[string]Session)
var data_b []Post1
var ID_G string
var User_G string
