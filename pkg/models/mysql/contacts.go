package mysql

import (
	"database/sql"

	"github.com/mbichoh/contactDash/pkg/models"
)

type ContactModel struct {
	DB *sql.DB
}

func (c *ContactModel) Insert(name string, contact string, uid int) (int, error) {

	stmt := `INSERT INTO contacts (name, contact, created_by_id) VALUES (?,?,?)`

	result, err := c.DB.Exec(stmt, name, contact, uid)
	if err != nil {
		return 0, nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

func (c *ContactModel) Get(id int) (*models.Contact, error) {
	stmt := `SELECT id, name, contact FROM contacts WHERE id = ?`

	row := c.DB.QueryRow(stmt, id)

	s := &models.Contact{}

	err := row.Scan(&s.ID, &s.Name, &s.Mobile)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

func (c *ContactModel) Latest(uid int) ([]*models.Contact, error) {

	stmt := `SELECT id, name, contact FROM contacts WHERE created_by_id = ? ORDER BY id DESC;`

	rows, err := c.DB.Query(stmt, uid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	contacts := []*models.Contact{}

	for rows.Next() {
		s := &models.Contact{}

		err = rows.Scan(&s.ID, &s.Name, &s.Mobile)
		if err != nil {
			return nil, err
		}

		contacts = append(contacts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return contacts, nil
}

func (c *ContactModel) Delete(id_no int) (int, error) {
	stmt := `DELETE FROM contacts WHERE ID = ?`
	result, err := c.DB.Exec(stmt, id_no)

	if err != nil {
		return 0, nil
	}
	id, err := result.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

func (c *ContactModel) Update(name, contact, idn string) (int, error) {

	stmt := `UPDATE contacts SET name = ?, contact = ? WHERE id = ?`

	result, err := c.DB.Exec(stmt, name, contact, idn)

	if err != nil {
		return 0, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (c *ContactModel) GetGroupedContacts(group_id int) ([]*models.Contact, error) {

	stmt := `select c1.id, c1.name, c1.contact from contacts
					c1 inner join groups c3 inner join grouped_contacts c2 
						on c1.id = c2.contact_id and c3.id = c2.group_id 
							where c2.group_id = ?;`

	rows, err := c.DB.Query(stmt, group_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	groups := []*models.Contact{}

	for rows.Next() {
		s := &models.Contact{}

		err = rows.Scan(&s.ID, &s.Name, &s.Mobile)
		if err != nil {
			return nil, err
		}

		groups = append(groups, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return groups, nil

}
