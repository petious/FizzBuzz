package sqlUtils

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"time"
)

type RequestFizzBuzz struct {
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Limit int    `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}

func UpdateDBWithRequest(db *sqlx.DB, int1 int, int2 int, limit int, str1 string, str2 string) error {
	alreadyExist, err := CheckIfRequestAlreadyInTable(db, int1, int2, limit, str1, str2)
	if err != nil {
		return err
	}

	if alreadyExist {
		err = UpdateCountRequest(db, int1, int2, limit, str1, str2)
		return err
	}
	err = AddRequestToDB(db, int1, int2, limit, str1, str2)
	return err
}

func AddRequestToDB(db *sqlx.DB, int1 int, int2 int, limit int, str1 string, str2 string) error {

	query := "INSERT INTO request_history (`int1`, `int2`, `limit`, `str1`, `str2`) VALUES (?,?,?,?,?) "

	_, err := db.Exec(query, int1, int2, limit, str1, str2)
	return err
}

func UpdateCountRequest(db *sqlx.DB, int1 int, int2 int, limit int, str1 string, str2 string) error {

	query := "UPDATE request_history SET `count`=`count`+1 WHERE `int1` = ? AND `int2` = ? AND `limit` = ? AND str1 = ? AND str2 = ? "

	_, err := db.Exec(query, int1, int2, limit, str1, str2)
	return err
}

func CheckIfRequestAlreadyInTable(db *sqlx.DB, int1 int, int2 int, limit int, str1 string, str2 string) (bool, error) {
	var alreadyExist bool
	query := "SELECT EXISTS(SELECT * FROM request_history WHERE `int1` = ? AND `int2` = ? AND `limit` = ? AND str1 = ? AND str2 = ?)"

	row := db.QueryRow(query, int1, int2, limit, str1, str2)
	if err := row.Scan(&alreadyExist); err != nil {
		return false, err
	}
	return alreadyExist, nil
}

func ConnectDBwithRetry(connectionString string, retryDurationSeconds int) (*sqlx.DB, error) {
	retryCount := 0

	var db *sqlx.DB
	var err error
	for retryCount < retryDurationSeconds {
		db, err = sqlx.Connect("mysql", connectionString)
		if err != nil {
			log.Warnf("Failed to connect to database on retry %d, err: %s", retryCount, err)
			retryCount += 1
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	if retryCount >= retryDurationSeconds {
		return nil, err
	}
	return db, nil
}

func GetDBConnectionString(user string, pwd string, host string, port string, db string) string {
	return user + ":" + pwd + "@(" + host + ":" + port + ")/" + db
}

func GetMostFrequentRequest(db *sqlx.DB) ([]RequestFizzBuzz, int, error) {
	mostFrequentRequests := make([]RequestFizzBuzz, 0)
	var count int
	query := "SELECT `int1`, `int2`, `limit`, `str1`, `str2`, `count` FROM request_history WHERE `count` = ( SELECT MAX(`count`) FROM request_history ) ;"
	rows, err := db.Query(query)
	if err != nil {
		log.Error("Error getting most frequent request, ", err)
		return mostFrequentRequests, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var mostFrequentRequest RequestFizzBuzz
		if err = rows.Scan(&mostFrequentRequest.Int1, &mostFrequentRequest.Int2, &mostFrequentRequest.Limit,
			&mostFrequentRequest.Str1, &mostFrequentRequest.Str2, &count); err != nil {
			return mostFrequentRequests, 0, err
		}
		mostFrequentRequests = append(mostFrequentRequests, mostFrequentRequest)
	}

	return mostFrequentRequests, count, nil
}

func GetRequestCount(db *sqlx.DB, int1 int, int2 int, limit int, str1 string, str2 string) (int, error) {
	var count int
	query := "SELECT `count` FROM request_history WHERE `int1` = ? AND `int2` = ? AND `limit` = ? AND str1 = ? AND str2 = ?"
	row := db.QueryRow(query, int1, int2, limit, str1, str2)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func TruncateTable(db *sqlx.DB) error {
	query := `TRUNCATE TABLE request_history;`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
