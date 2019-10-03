package mysql

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"

	"github.com/mbichoh/contactDash/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name string, mobile string, password string, token int, verified bool) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO contact_users (name, mobile, hashed_password, token, isVerified) VALUES(?, ?, ?, ?, ?)`
	stmtCheck := `SELECT * FROM contact_users WHERE mobile = ?`
	row := u.DB.QueryRow(stmtCheck, mobile)
	if row.Scan(&mobile) != sql.ErrNoRows {
		return models.ErrDuplicateNumber
	}
	_, err = u.DB.Exec(stmt, name, mobile, string(hashedPassword), token, verified)
	if err != nil {
		return err
	}

	return err

}

func (u *UserModel) Authenticate(mobile string, password string) (int, error) {

	var id int
	var hashedPassword []byte
	var isVer bool = true

	row := u.DB.QueryRow("SELECT id, hashed_password  FROM contact_users WHERE mobile = ? AND isVerified = true ", mobile)
	err := row.Scan(&id, &hashedPassword)

	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	if isVer != true {
		return 0, models.ErrNotVerified
	}

	return id, nil
}

func (u *UserModel) Get(id int) (*models.User, error) {
	s := &models.User{}
	stmt := `SELECT id, name, mobile FROM contact_users WHERE id = ?`
	err := u.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Name, &s.Mobile)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

func (u *UserModel) Verify(token int) (*models.User, error) {

	stmt := `SELECT id, mobile FROM contact_users WHERE token = ? AND isVerified = false`
	row := u.DB.QueryRow(stmt, token)

	s := &models.User{}

	if row.Scan(&s.Mobile) == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	}

	return s, nil
}

func (u *UserModel) IsVerified(idn int) (int, error) {
	stmtCheck := `UPDATE contact_users SET token = 0, isVerified = true  WHERE token = ?`
	result, err := u.DB.Exec(stmtCheck, idn)

	if err != nil {
		return 0, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
