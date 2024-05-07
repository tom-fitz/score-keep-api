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

func TestImportLeagues(t *testing.T) {
	teamsFile, err := os.CreateTemp("", "teams.csv")
	assert.NoError(t, err)
	defer os.Remove(teamsFile.Name())
	teamsFile.WriteString("Team1,Captain1,2021\nTeam2,Captain2,2022\n")
	teamsFile.Close()

	playersFile, err := os.CreateTemp("", "players.csv")
	assert.NoError(t, err)
	defer os.Remove(playersFile.Name())
	playersFile.WriteString("John,Doe,john@example.com,1234567890,123,A,Team1\nJane,Smith,jane@example.com,9876543210,456,B,Team2\n")
	playersFile.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", "1")
	writer.CreateFormFile("teams", teamsFile.Name())
	writer.CreateFormFile("players", playersFile.Name())
	writer.Close()

	req, err := http.NewRequest("POST", "/leagues/1/import", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM score_keep_db.public.leagues WHERE id = \\$1").
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{
		{Key: "id", Value: "1"},
	}

	h := &Handler{db: db}

	h.importLeagues(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Import successful", response["message"])

	assert.NoError(t, mock.ExpectationsWereMet())
}
