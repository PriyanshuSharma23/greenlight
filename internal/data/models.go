package data

import (
	"database/sql"
	"errors"
)

var (
	ErrNoRecordFound  = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("duplicate email")
)

type Models struct {
	Movies      MovieModel
	Permissions PermissionModel
	Users       UserModel
	Tokens      TokensModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		MovieModel{DB: db},
		PermissionModel{DB: db},
		UserModel{DB: db},
		TokensModel{DB: db},
	}
}
