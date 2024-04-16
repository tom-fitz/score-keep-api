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
		email := p.email
		firstName := p.firstName
		lastName := p.lastName
		level := p.level
		phone := p.phone
		teamNames := p.teamNames
		usaNum := p.usaNum

		emailQuery := "SELECT COUNT(*) FROM score_keep_db.public.players WHERE email = $1"
		var existingPlayer *Player
		err := db.QueryRow(emailQuery, p.email).Scan(&existingPlayer)
		if err != nil {
			return fmt.Errorf("error finding player count: %w", err)
		}

		if existingPlayer != nil {
			if existingPlayer.usaNum != p.usaNum {
				updatePlayerQuery := `
					UPDATE score_keep_db.public.players
					SET usaNum = $1, firstName = $2, lastName = $3, level = $4, phone = $5
					WHERE email = $6
				`
				_, err = db.Exec(updatePlayerQuery, usaNum, firstName, lastName, level, phone, email)
				if err != nil {
					return fmt.Errorf("error updating player: %w", err)
				}

				removeOldUSANumberQuery := `
					DELETE FROM score_keep_db.public.player_team
					WHERE usaNum = $1
				`
				_, err = db.Exec(removeOldUSANumberQuery, existingPlayer.usaNum)
				if err != nil {
					return fmt.Errorf("error removing old USA number: %w", err)
				}

				insertNewUSANumberQuery := `
					INSERT INTO score_keep_db.public.player_team (usaNum, teamName)
					VALUES ($1, $2)
				`
				for _, team := range teamNames {
					_, err := db.Exec(insertNewUSANumberQuery, p.usaNum, team)
					if err != nil {
						return fmt.Errorf("error inserting new USA number with team: %w", err)
					}
				}
			}
		}

		teams := strings.Split(teamNames, ", ")

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin transaction: %w", err)
		}

		query := "INSERT INTO score_keep_db.public.players (email, firstName, lastName, level, phone, usaNum) VALUES ($1, $2, $3, $4, $5, $6)"
		_, err = tx.Exec(query, email, firstName, lastName, level, phone, usaNum)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("insert player data: %w", err)
		}

		for _, team := range teams {
			query = "INSERT INTO score_keep_db.public.player_team (usaNum, teamName) VALUES ($1, $2)"
			_, err = tx.Exec(query, usaNum, team)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("insert player-team relationship: %w", err)
			}
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("commit transaction: %w", err)
		}
	}
	return nil
}

// InsertTeamData inserts team data into the database
func InsertTeamData(db *sql.DB, teams []map[string]string) error {
	for _, team := range teams {
		name := team["name"]
		captain := team["captain"]
		firstYear := team["firstYear"]

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
