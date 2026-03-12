package repository

import (
	"database/sql"
)

type SearchRepository interface {
	AdvancedSearch(query string, boardID *int, userID int) ([]map[string]interface{}, error)
}

type searchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) SearchRepository {
	return &searchRepository{db: db}
}

func (r *searchRepository) AdvancedSearch(searchQuery string, boardID *int, userID int) ([]map[string]interface{}, error) {
	q := `
		SELECT c.id, c.title, c.description, c.label_color, l.title as list_title, b.title as board_title
		FROM cards c
		JOIN lists l ON c.list_id = l.id
		JOIN boards b ON l.board_id = b.id
		LEFT JOIN board_collaborators bc ON b.id = bc.board_id
		WHERE (b.user_id = ? OR bc.user_id = ?)
		AND (c.title LIKE ? OR c.description LIKE ?)
	`
	
	args := []interface{}{userID, userID, "%" + searchQuery + "%", "%" + searchQuery + "%"}
	
	if boardID != nil {
		q += " AND b.id = ?"
		args = append(args, *boardID)
	}
	
	q += " LIMIT 20"
	
	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var title, listTitle, boardTitle string
		var desc, labelColor sql.NullString
		
		if err := rows.Scan(&id, &title, &desc, &labelColor, &listTitle, &boardTitle); err != nil {
			return nil, err
		}
		
		res := map[string]interface{}{
			"id":          id,
			"title":       title,
			"list_name":   listTitle,
			"board_name":  boardTitle,
		}
		if desc.Valid { res["description"] = desc.String }
		if labelColor.Valid { res["label_color"] = labelColor.String }
		
		results = append(results, res)
	}
	return results, nil
}
