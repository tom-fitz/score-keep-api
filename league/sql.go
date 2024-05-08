package league

import (
	"database/sql"
	"fmt"
)

func validateLeague(id string, db *sql.DB) (League, error) {
	query := "SELECT * FROM score_keep_db.public.leagues WHERE id = $1"
	var league League
	err := db.QueryRow(query, id).Scan(&league.Id, &league.Name, &league.Level)
	if err != nil {
		return league, fmt.Errorf("error finding league: %w", err)
	}
	if &league == nil {
		return league, fmt.Errorf("league not found with id %s", id)
	}
	return league, nil
}
