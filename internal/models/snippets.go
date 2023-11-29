package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// Insert will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	// SQL statement for inserting a new record in the snippets table, and storing
	// the resulting ID.

	stmt := `INSERT INTO snippets (title, content, created, expires)
			VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	// Use the Exec() method on the embedded connection pool to execute the
	// query, with the placeholder values.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	// Use the LastInsertId() method on the result object to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Get will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
             WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// Returns a sql.Row pointer which holds the result from the database.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct.
	var s Snippet
	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRows
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

// Latest will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	// Returns a sql.Rows pointer which acts as an iterator over the returned
	// records.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// We defer rows.Close() to ensure the sql.Rows resultset is always properly
	// closed before the Latest() method returns. This defer statement should come
	// *after* you check for an error from the Query() method. Otherwise, if Query()
	// returns an error, you'll get a panic trying to close a nil resultset.
	defer rows.Close()
	// Initialize an empty slice to hold the models.Snippets objects.
	var snippets []Snippet
	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		var s Snippet
		// Use rows.Scan() to copy the values from each field in the row
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// Call rows.Err() to retrieve any error that was encountered during iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK, then return the Snippets slice.
	return snippets, nil
}
