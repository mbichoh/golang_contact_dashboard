package mysql

import (
	"database/sql"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-sql-driver/mysql"
	"github.com/mbichoh/contactDash/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, mobile, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO contact_users (name, mobile, hashed_password) VALUES(?, ?, ?)`

	_, err = u.DB.Exec(stmt, name, mobile, string(hashedPassword))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "users_uc_mobile") {
				return models.ErrDuplicateNumber
			}
			// CHECK: what if mysql change their error number and error messages? Your application will break
		}
	}
	return err

}

func (u *UserModel) Authenticate(mobile, password string) (int, error) {

	var id int
	var hashedPassword []byte

	row := u.DB.QueryRow("SELECT id, hashed_password FROM contact_users WHERE mobile = ?", mobile)
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
