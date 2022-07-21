package main

import (
	"Fizzbuzz/sqlUtils"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

const DbConnectionString = "root:if5!spl?@(172.21.0.2:3306)/fizzbuzz_test"
const RetryDurationSeconds = 60

func TestFizzBuzzApiResponseCode(t *testing.T) {
	var app App
	var err error
	app.DB, err = sqlUtils.ConnectDBwithRetry(DbConnectionString, RetryDurationSeconds)
	if err != nil {
		t.Error("failed to connect to database", err)
		return
	}

	router := app.SetUpRouter()
	reqWithoutQueryParam, _ := http.NewRequest("GET", "/fizzbuzz", nil)
	response := executeRequest(reqWithoutQueryParam, router)
	checkResponseCode(t, response.Code, 400)

	reqWithIncorrectQueryParam, _ := http.NewRequest("GET", "/fizzbuzz?int1=1&int2=pasunint&limit=3&str1=fizz&str2=buzz", nil)
	response = executeRequest(reqWithIncorrectQueryParam, router)
	checkResponseCode(t, response.Code, 400)

	reqWithCorrectQueryParam, _ := http.NewRequest("GET", "/fizzbuzz?int1=1&int2=2&limit=3&str1=fizz&str2=buzz", nil)
	response = executeRequest(reqWithCorrectQueryParam, router)
	checkResponseCode(t, response.Code, 200)

	if err = sqlUtils.TruncateTable(app.DB); err != nil {
		t.Error("Error truncating table")
	}

}

func executeRequest(req *http.Request, router *gin.Engine) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
