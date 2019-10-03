package mysql

import (
	"database/sql"

	"github.com/mbichoh/contactDash/pkg/models"
)

type GroupsModel struct {
	DB *sql.DB
}

func (g *GroupsModel) GroupInsertName(name string) (int, error) {

	stmt := `INSERT INTO groups(group_name) VALUES (?)`

	result, err := g.DB.Exec(stmt, name)
	if err != nil {
		return 0, nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}
func (g *GroupsModel) GroupFetchNames(id int) ([]*models.Groups, error) {

	stmt := `select distinct c3.id, c3.group_name from contacts
		 		c1 inner join groups c3 inner join grouped_contacts c2 
				 on c1.id = c2.contact_id and c3.id = c2.group_id 
					 where c1.created_by_id = ?;`

	rows, err := g.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	groups := []*models.Groups{}

	for rows.Next() {
		s := &models.Groups{}

		err = rows.Scan(&s.ID, &s.Name)
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

func (g *GroupsModel) Get(id int) (*models.Groups, error) {
	stmt := `SELECT id, group_name FROM groups WHERE id = ?`

	row := g.DB.QueryRow(stmt, id)

	s := &models.Groups{}

	err := row.Scan(&s.ID, &s.Name)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, nil
}
