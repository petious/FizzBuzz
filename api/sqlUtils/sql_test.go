package sqlUtils

import (
	"github.com/go-playground/assert/v2"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

const DbConnectionString = "root:if5!spl?@(172.21.0.2:3306)/fizzbuzz_test"
const RetryDurationSeconds = 60

func TestAddRequestToDB(t *testing.T) {
	db, err := ConnectDBwithRetry(DbConnectionString, RetryDurationSeconds)
	if err != nil {
		t.Error("failed to connect to database", err)
		return
	}

	if err = TruncateTable(db); err != nil {
		t.Error("Error truncating table")
		return
	}

	if err = AddRequestToDB(db, 1, 2, 3, "fizz", "buzz"); err != nil {
		t.Error("Error adding request to database", err)
		return
	}

	alreadyExist, err := CheckIfRequestAlreadyInTable(db, 1, 2, 3, "fizz", "buzz")
	if err != nil {
		t.Error("Error checking if request already exists in table,", err)
	}

	if !alreadyExist {
		t.Error("Error adding request to database")
	}

	if err = UpdateCountRequest(db, 1, 2, 3, "fizz", "buzz"); err != nil {
		t.Error("Error updating request count", err)
	}

	count, err := GetRequestCount(db, 1, 2, 3, "fizz", "buzz")
	if err != nil {
		t.Error("Error getting request count", err)
	}
	if count != 2 {
		t.Error("Request count is not properly updated")
	}

	mostFrequentRequests, count, err := GetMostFrequentRequest(db)
	if err != nil {
		t.Error("Error getting most frequent request,", err)
	}

	assert.Equal(t, mostFrequentRequests, []RequestFizzBuzz{{1, 2, 3, "fizz", "buzz"}})

	if err = TruncateTable(db); err != nil {
		t.Error("Error truncating table")
	}
}
