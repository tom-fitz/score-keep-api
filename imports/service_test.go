package imports

import (
	"bytes"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func setupTestService() (*Service, sqlmock.Sqlmock, *httptest.ResponseRecorder, *gin.Context, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h := &Service{db: db}
	return h, mock, w, c, nil
}

func setupTestRouter(method string, url string, body *bytes.Buffer) *http.Request {
	req, _ := http.NewRequest(method, url, body)
	return req
}
