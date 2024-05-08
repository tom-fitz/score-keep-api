package imports

import (
	"database/sql"
	"fmt"
	"strings"
)

func validateLeague(id string, db *sql.DB) error {
	query := "SELECT COUNT(*) FROM score_keep_db.public.leagues WHERE id = $1"
	var count int
	err := db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("error finding league: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("league not found with id %s", id)
	}
	return nil
}

// InsertPlayerData inserts player data into the database
func InsertPlayerData(db *sql.DB, players []Player) error {
	for _, p := range players {
		email := p.Email
		firstName := p.FirstName
		lastName := p.LastName
		level := p.Level
		phone := p.Phone
		teamNames := p.TeamNames
		usaNum := p.Usanum

		//emailQuery := "SELECT COUNT(*) FROM score_keep_db.public.players WHERE email = $1"
		usaNumQuery := "SELECT COUNT(*) FROM score_keep_db.public.players WHERE usanum = $1"
		var count int
		err := db.QueryRow(usaNumQuery, p.Email).Scan(&count)
		if err != nil {
			return fmt.Errorf("error checking existing player: %w", err)
		}

		if count > 0 {
			// Player already exists, update their details
			updatePlayerQuery := `
                UPDATE score_keep_db.public.players
                SET usanum = $1, firstName = $2, lastName = $3, level = $4, phone = $5
                WHERE email = $6
            `
			_, err := db.Exec(updatePlayerQuery, usaNum, firstName, lastName, level, phone, email)
			if err != nil {
				return fmt.Errorf("error updating player: %w", err)
			}

			// Remove old player-team relationships
			removeOldTeamQuery := `
                DELETE FROM score_keep_db.public.player_team
                WHERE usanum = $1
            `
			_, err = db.Exec(removeOldTeamQuery, usaNum)
			if err != nil {
				return fmt.Errorf("error removing old player-team relationships: %w", err)
			}
		} else {
			// Player doesn't exist, insert a new one
			insertPlayerQuery := `
                INSERT INTO score_keep_db.public.players (email, firstName, lastName, level, phone, usanum)
                VALUES ($1, $2, $3, $4, $5, $6)
            `
			_, err := db.Exec(insertPlayerQuery, email, firstName, lastName, level, phone, usaNum)
			if err != nil {
				return fmt.Errorf("error inserting new player: %w", err)
			}
		}

		// Insert player-team relationships
		teams := strings.Split(teamNames, ", ")
		insertTeamQuery := `
            INSERT INTO score_keep_db.public.player_team (usanum, teamName)
            VALUES ($1, $2)
        `
		for _, team := range teams {
			_, err := db.Exec(insertTeamQuery, usaNum, team)
			if err != nil {
				return fmt.Errorf("error inserting player-team relationship: %w", err)
			}
		}
	}
	return nil
}

// InsertTeamData inserts team data into the database
func InsertTeamData(db *sql.DB, teams []Team) error {
	for _, team := range teams {
		name := team.Name
		captain := team.Captain
		firstYear := team.FirstYear

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin transaction: %w", err)
		}

		query := "INSERT INTO score_keep_db.public.teams (name, captain, firstYear) VALUES ($1, $2, $3)"
		_, err = tx.Exec(query, name, captain, firstYear)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("insert team data: %w", err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("commit transaction: %w", err)
		}
	}

	return nil
}
