package forum

import "database/sql"

var db *sql.DB
var sessionsMap = make(map[string]Session)
var data_b []Post1
var ID_G string
var User_G string
var (
    googleClientID     = "YOUR_GOOGLE_CLIENT_ID"
    googleClientSecret = "YOUR_GOOGLE_CLIENT_SECRET"
    githubClientID     = "YOUR_GITHUB_CLIENT_ID"
    githubClientSecret = "YOUR_GITHUB_CLIENT_SECRET"
)
