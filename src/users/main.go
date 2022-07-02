package users

import (
	"database/sql"
	"math/rand"
)

var DB *sql.DB
var R *rand.Rand
var LocalHost string
