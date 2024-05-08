package imports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func createTestTeamFile(t *testing.T) string {
	teamsFile, err := os.CreateTemp("", "teams.csv")
	assert.NoError(t, err)
	teamsFile.WriteString("Team1,Captain1,2021\nTeam2,Captain2,2022\n")
	teamsFile.Close()

	return teamsFile.Name()
}

func TestImportTeams(t *testing.T) {
	teamsFilename := createTestTeamFile(t)
	defer func() {
		os.Remove(teamsFilename)
	}()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", "1")
	writer.CreateFormFile("teams", teamsFilename)
	writer.Close()

	t.Run("League Id not found", func(t *testing.T) {
		h, _, w, c, err := setupTestService()
		assert.NoError(t, err)
		defer h.db.Close()

		req := setupTestRouter("POST", "", body)
		c.Request = req

		h.ImportTeams(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "league id required")
	})

	t.Run("League not found", func(t *testing.T) {
		h, mock, w, c, err := setupTestService()
		assert.NoError(t, err)
		defer h.db.Close()

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM score_keep_db.public.leagues WHERE id = \\$1").
			WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		req := setupTestRouter("POST", "/v1/league/1/teams/import", body)

		c.Request = req
		c.Params = gin.Params{
			{Key: "leagueId", Value: "1"},
		}

		h.ImportTeams(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "league not found")
	})

	t.Run("Teams file read error", func(t *testing.T) {
		h, mock, w, c, err := setupTestService()
		assert.NoError(t, err)
		defer h.db.Close()

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM score_keep_db.public.leagues WHERE id = \\$1").
			WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		body := new(bytes.Buffer)
		req := setupTestRouter("POST", "/v1/league/1/teams/import", body)
		c.Request = req
		c.Params = gin.Params{
			{Key: "leagueId", Value: "1"},
		}

		h.ImportTeams(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "teams.csv file not provided")
	})

	t.Run("InsertTeamData error", func(t *testing.T) {
		h, mock, w, c, err := setupTestService()
		assert.NoError(t, err)
		defer h.db.Close()

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM score_keep_db.public.leagues WHERE id = \\$1").
			WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO score_keep_db.public.teams \\(name, captain, firstYear\\) VALUES \\(\\$1, \\$2, \\$3\\)").
			WithArgs("Team1", "Captain1", "2021").
			WillReturnError(fmt.Errorf("insert team data error"))
		mock.ExpectCommit()

		teamsFile, err := os.Open(teamsFilename)
		assert.NoError(t, err)
		defer teamsFile.Close()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("teams", teamsFilename)
		assert.NoError(t, err)
		_, err = io.Copy(part, teamsFile)
		assert.NoError(t, err)
		writer.Close()

		req := setupTestRouter("POST", "/v1/league/1/teams/import", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		c.Request = req
		c.Params = gin.Params{
			{Key: "leagueId", Value: "1"},
		}

		h.ImportTeams(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "insert team data")
	})
}
