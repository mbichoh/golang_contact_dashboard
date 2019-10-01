package mysql

import (
	"database/sql"
)

type GroupedContactsModel struct {
	DB *sql.DB
}

func (gc *GroupedContactsModel) Insert(contactid, groupnameid int) (int, error) {

	stmt := `INSERT INTO grouped_contacts (contact_id, group_id) VALUES (?,?)`

	result, err := gc.DB.Exec(stmt, contactid, groupnameid)
	if err != nil {
		return 0, nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

func (gc *GroupedContactsModel) DeleteContact(id_no int) (int, error) {
	stmt := `delete from grouped_contacts where contact_id = ?;`
	result, err := gc.DB.Exec(stmt, id_no)

	if err != nil {
		return 0, nil
	}
	id, err := result.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return int(id), nil

}

func (gc *GroupedContactsModel) DeleteGroup(idno int) (int, error) {

	stmt := `delete from groups where id = ?;`

	result, err := gc.DB.Exec(stmt, idno)

	if err != nil {
		return 0, nil
	}
	id, err := result.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}
