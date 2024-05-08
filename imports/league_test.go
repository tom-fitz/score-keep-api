package imports

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupTestService() (*Handler, sqlmock.Sqlmock, *httptest.ResponseRecorder, *gin.Context, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h := &Handler{db: db}
	return h, mock, w, c, nil
}

func setupTestRouter(method string, url string, body *bytes.Buffer) *http.Request {
	req, _ := http.NewRequest(method, url, body)
	return req
}

func createTestFiles(t *testing.T) (string, string) {
	teamsFile, err := os.CreateTemp("", "teams.csv")
	assert.NoError(t, err)
	teamsFile.WriteString("Team1,Captain1,2021\nTeam2,Captain2,2022\n")
	teamsFile.Close()

	playersFile, err := os.CreateTemp("", "players.csv")
	assert.NoError(t, err)
	playersFile.WriteString("John,Doe,john@example.com,1234567890,123,A,Team1\nJane,Smith,jane@example.com,9876543210,456,B,Team2\n")
	playersFile.Close()

	return teamsFile.Name(), playersFile.Name()
}

func TestImportLeagues(t *testing.T) {
	teamsFilename, playersFilename := createTestFiles(t)
	defer func() {
		os.Remove(teamsFilename)
		os.Remove(playersFilename)
	}()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", "1")
	writer.CreateFormFile("teams", teamsFilename)
	writer.CreateFormFile("players", playersFilename)
	writer.Close()

	t.Run("League Id not found", func(t *testing.T) {
		h, _, w, c, err := setupTestService()
		assert.NoError(t, err)
		defer h.db.Close()

		req := setupTestRouter("POST", "", body)
		c.Request = req

		h.importLeagues(c)
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

		req := setupTestRouter("POST", "/league/1/import", body)

		c.Request = req
		c.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		h.importLeagues(c)

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
		writer := multipart.NewWriter(body)
		writer.WriteField("id", "1")
		writer.CreateFormFile("teams", "nonexistent_file.csv")
		writer.CreateFormFile("players", playersFilename)
		writer.Close()

		req := setupTestRouter("POST", "/leagues/1/import", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		c.Request = req
		c.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		h.importLeagues(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "teams.csv file not provided")
	})

	t.Run("Import successful", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/leagues/1/import", body)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		h, mock, w, c, err := setupTestService()
		assert.NoError(t, err)
		defer h.db.Close()

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM score_keep_db.public.leagues WHERE id = \\$1").
			WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		c.Request = req
		c.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		h.importLeagues(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Import successful", response["message"])
	})
}
