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

func createTestPlayerFile(t *testing.T) string {
	playersFile, err := os.CreateTemp("", "players.csv")
	assert.NoError(t, err)
	playersFile.WriteString("John,Doe,john@example.com,1234567890,123,A,Team1\nJane,Smith,jane@example.com,9876543210,456,B,Team2\n")
	playersFile.Close()

	return playersFile.Name()
}

func TestImportPlayers(t *testing.T) {
	playersFilename := createTestPlayerFile(t)
	defer func() {
		os.Remove(playersFilename)
	}()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", "1")
	writer.CreateFormFile("players", playersFilename)
	writer.Close()

	t.Run("League Id not found", func(t *testing.T) {
		h, _, w, c, err := setupTestService()
		assert.NoError(t, err)
		defer h.db.Close()

		req := setupTestRouter("POST", "", body)
		c.Request = req

		h.ImportPlayers(c)
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

		req := setupTestRouter("POST", "/v1/league/1/players/import", body)

		c.Request = req
		c.Params = gin.Params{
			{Key: "leagueId", Value: "1"},
		}

		h.ImportPlayers(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "league not found")
	})

	t.Run("Players file read error", func(t *testing.T) {
		h, mock, w, c, err := setupTestService()
		assert.NoError(t, err)
		defer h.db.Close()

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM score_keep_db.public.leagues WHERE id = \\$1").
			WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		body := new(bytes.Buffer)
		req := setupTestRouter("POST", "/v1/league/1/players/import", body)
		c.Request = req
		c.Params = gin.Params{
			{Key: "leagueId", Value: "1"},
		}

		h.ImportPlayers(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "players.csv file not provided")
	})

	t.Run("InsertPlayerData error", func(t *testing.T) {
		h, mock, w, c, err := setupTestService()
		assert.NoError(t, err)
		defer h.db.Close()

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM score_keep_db.public.leagues WHERE id = \\$1").
			WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO score_keep_db.public.players \\(firstName, lastName, email, phone, usaNum, level, teamNames\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7\\)").
			WithArgs("test", "name", "test@test.com", "111111111", "098SADF089", "3", "test, one, two").
			WillReturnError(fmt.Errorf("insert player data error"))
		mock.ExpectCommit()

		playersFile, err := os.Open(playersFilename)
		assert.NoError(t, err)
		defer playersFile.Close()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("players", playersFilename)
		assert.NoError(t, err)
		_, err = io.Copy(part, playersFile)
		assert.NoError(t, err)
		writer.Close()

		req := setupTestRouter("POST", "/v1/league/1/players/import", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		c.Request = req
		c.Params = gin.Params{
			{Key: "leagueId", Value: "1"},
		}

		h.ImportPlayers(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "insert player data")
	})
}
