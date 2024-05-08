package imports

import (
	"database/sql"
	"fmt"
	"strings"
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

		usaNumQuery := "SELECT COUNT(*) FROM score_keep_db.public.players WHERE usanum = $1"
		var count int
		err := db.QueryRow(usaNumQuery, usaNum).Scan(&count)
		if err != nil {
			return fmt.Errorf("error checking existing player: %w", err)
		}

		if count > 0 {
			// Player already exists, update their details
			updatePlayerQuery := `
                UPDATE score_keep_db.public.players
                SET firstName = $1, lastName = $2, level = $3, phone = $4
                WHERE usanum = $5
            `
			_, err := db.Exec(updatePlayerQuery, firstName, lastName, level, phone, usaNum)
			if err != nil {
				return fmt.Errorf("error updating player: %w", err)
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

		// Remove old player-team relationships
		removeOldTeamQuery := `
            DELETE FROM score_keep_db.public.player_team
            WHERE usanum = $1
        `
		_, err = db.Exec(removeOldTeamQuery, usaNum)
		if err != nil {
			return fmt.Errorf("error removing old player-team relationships: %w", err)
		}

		// Insert player-team relationships
		teams := strings.Split(teamNames, ", ")
		insertTeamQuery := `
     INSERT INTO score_keep_db.public.player_team (usanum, team_name, team_id)
     VALUES($1, $2, $3)
     ON CONFLICT (usanum, team_id, team_name) DO NOTHING
`
		for _, team := range teams {
			var teamId int
			getTeamIdQuery := `SELECT id FROM teams WHERE name = $1`
			err := db.QueryRow(getTeamIdQuery, team).Scan(&teamId)
			if err != nil {
				if err == sql.ErrNoRows {
					// Team doesn't exist in the teams table, skip inserting the relationship
					continue
				}
				return fmt.Errorf("error finding team id from teams table: %w", err)
			}

			_, err = db.Exec(insertTeamQuery, usaNum, team, teamId)
			if err != nil {
				return fmt.Errorf("error inserting player-team relationship: %w", err)
			}
		} //		// Insert player-team relationships
		//		teams := strings.Split(teamNames, ", ")
		//		insertTeamQuery := `
		//     INSERT INTO score_keep_db.public.player_team (usanum, team_name, team_id)
		//     VALUES($1, $2, $3)
		//`
		//		for _, team := range teams {
		//			var teamId int
		//			getTeamIdQuery := `SELECT id FROM teams WHERE name = $1`
		//			err := db.QueryRow(getTeamIdQuery, team).Scan(&teamId)
		//			if err != nil {
		//				if err == sql.ErrNoRows {
		//					// Team doesn't exist in the teams table, skip inserting the relationship
		//					continue
		//				}
		//				return fmt.Errorf("error finding team id from teams table: %w", err)
		//			}
		//
		//			// Check if the player-team relationship already exists
		//			checkRelationshipQuery := `
		//        SELECT COUNT(*) FROM score_keep_db.public.player_team
		//        WHERE usanum = $1 AND team_id = $2
		//   `
		//			var count int
		//			err = db.QueryRow(checkRelationshipQuery, usaNum, teamId).Scan(&count)
		//			if err != nil {
		//				return fmt.Errorf("error checking player-team relationship: %w", err)
		//			}
		//
		//			if count > 0 {
		//				// Player-team relationship already exists, skip insertion
		//				continue
		//			}
		//
		//			_, err = db.Exec(insertTeamQuery, usaNum, team, teamId)
		//			if err != nil {
		//				return fmt.Errorf("error inserting player-team relationship: %w", err)
		//			}
		//		}
	}
	return nil
}

// InsertTeamData inserts team data into the database
func InsertTeamData(db *sql.DB, teams []Team, league League) error {
	for _, team := range teams {
		name := team.Name
		captain := team.Captain
		firstYear := team.FirstYear

		// Check if the team already exists
		teamExistsQuery := "SELECT id FROM score_keep_db.public.teams WHERE name = $1"
		var existingTeamID int
		err := db.QueryRow(teamExistsQuery, name).Scan(&existingTeamID)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("error checking existing team: %w", err)
		}

		var teamID int
		if existingTeamID != 0 {
			// Team already exists, update their details
			updateTeamQuery := `
                UPDATE score_keep_db.public.teams
                SET captain = $1, firstYear = $2
                WHERE name = $3
                RETURNING id
            `
			err := db.QueryRow(updateTeamQuery, captain, firstYear, name).Scan(&teamID)
			if err != nil {
				return fmt.Errorf("error updating team: %w", err)
			}
		} else {
			// Team doesn't exist, insert a new one
			insertTeamQuery := `
                INSERT INTO score_keep_db.public.teams (name, captain, firstYear)
                VALUES ($1, $2, $3)
                RETURNING id
            `
			err := db.QueryRow(insertTeamQuery, name, captain, firstYear).Scan(&teamID)
			if err != nil {
				return fmt.Errorf("error inserting new team: %w", err)
			}
		}

		// Insert or update the league-team relationship
		insertLeagueTeamQuery := `
            INSERT INTO score_keep_db.public.league_team (league_id, league_name, team_name, team_id)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (league_id, team_id) DO UPDATE SET league_name = $2, team_name = $3
        `
		_, err = db.Exec(insertLeagueTeamQuery, league.Id, league.Name, name, teamID)
		if err != nil {
			return fmt.Errorf("error inserting league-team relationship: %w", err)
		}
	}

	return nil
}
