package models

import "errors"

var ErrNoRecord = errors.New("models: no such record found")
var ErrInvalidCredentials = errors.New("models: invalid credentials")
var ErrDuplicateNumber = errors.New("models: number exists")
var ErrNotVerified = errors.New("models: Signed in, not verified")

type (
	User struct {
		ID             int
		Name, Mobile   string
		HashedPassword []byte
		Token          int
		IsVerified     string
	}

	Contact struct {
		ID           int
		Name, Mobile string
	}

	Groups struct {
		ID   int
		Name string
	}

	GroupedContacts struct {
		ID        int
		ContactID int
		GroupID   int
	}
)
