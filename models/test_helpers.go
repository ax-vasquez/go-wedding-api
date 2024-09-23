package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Convenience variable used to perform "bad" lookup tests (e.g., ensure querying non-existent data returns an empty result)
var NilUuid = "00000000-0000-0000-0000-000000000000"

// Convenience variable to keep easy reference to the UUID of the first user in the test data set ("Rupinder McNiel")
var FirstUserIdStr = "0ad1d80a-329b-4ffe-89c1-87af4d945953"

// Convenience variable to keep easy reference to the UUID of the first user's invitee in the test data set ("Suman Sousa")
var FirstUserInviteeIdStr = "007170d7-5633-4a44-9326-ddf9dce5a6ef"

// Convenience variable to keep easy reference to the UUID of the entree in the test data set ("Caprese pasta")
var FirstEntreeIdStr = "f8cd5ea3-bb29-42fc-9984-a6c37d8b99c3"

// Convenience variable to keep easy reference to the UUID of the hors doeuvres in the test data set ("Crab puff")
var FirstHorsDoeuvresIdStr = "3baf970f-1670-4b42-ba81-63168a2f21b8"

// All test users have the same password
var TestUserPassword = "ASdf12#$"

func loadTestUsers(c context.Context) error {
	users := []User{}
	userInvitees := []UserInvitee{}
	admins := []User{}
	userDataFile, err := os.ReadFile("../test-fixtures/users.json")
	if err != nil {
		return errors.New("There was a problem loading test user data from ./test-fixtures/users.json: " + err.Error())
	}
	userInviteeDataFile, err := os.ReadFile("./../test-fixtures/invitees.json")
	if err != nil {
		return errors.New("There was a problem loading test user data from ./test-fixtures/invitees.json: " + err.Error())
	}
	adminsDataFile, err := os.ReadFile("./../test-fixtures/admins.json")
	if err != nil {
		return errors.New("There was a problem loading test user data from ./test-fixtures/admins.json: " + err.Error())
	}
	err = json.Unmarshal(userDataFile, &users)
	if err != nil {
		return errors.New("There was a problem unmarshaling the JSON from file ./test-fixtures/users.json: " + err.Error())
	}
	err = json.Unmarshal(userInviteeDataFile, &userInvitees)
	if err != nil {
		return errors.New("There was a problem unmarshaling the JSON from file ./test-fixtures/invitees.json: " + err.Error())
	}
	err = json.Unmarshal(userInviteeDataFile, &userInvitees)
	if err != nil {
		return errors.New("There was a problem unmarshaling the JSON from file ./test-fixtures/invitees.json: " + err.Error())
	}
	err = json.Unmarshal(adminsDataFile, &admins)
	if err != nil {
		return errors.New("There was a problem unmarshaling the JSON from file ./test-fixtures/invitees.json: " + err.Error())
	}
	err = CreateUsers(c, &users)
	if err != nil {
		return errors.New("There was a problem creating test user records: " + err.Error())
	}
	err = CreateUserInvitees(c, &userInvitees)
	if err != nil {
		return errors.New("There was a problem creating test user invitee records: " + err.Error())
	}
	err = CreateUsers(c, &admins)
	if err != nil {
		return errors.New("There was a problem creating test user invitee records: " + err.Error())
	}
	return nil
}

func loadTestEntrees(c context.Context) error {
	records := []Entree{}
	recordsFile, err := os.ReadFile("../test-fixtures/entrees.json")
	if err != nil {
		return errors.New("There was a problem loading test user data from ./test-fixtures/entrees.json: " + err.Error())
	}
	err = json.Unmarshal(recordsFile, &records)
	if err != nil {
		return errors.New("There was a problem unmarshaling the JSON from file ./test-fixtures/entrees.json: " + err.Error())
	}
	err = CreateEntrees(c, &records)
	if err != nil {
		return errors.New("There was a problem creating the test entree records: " + err.Error())
	}
	return nil
}

func loadTestHorsDoeuvres(c context.Context) error {
	records := []HorsDoeuvres{}
	recordsFile, err := os.ReadFile("../test-fixtures/hors_doeuvres.json")
	if err != nil {
		return errors.New("There was a problem loading test user data from ./test-fixtures/hors_doeuvres.json: " + err.Error())
	}
	err = json.Unmarshal(recordsFile, &records)
	if err != nil {
		return errors.New("There was a problem unmarshaling the JSON from file ./test-fixtures/hors_doeuvres.json: " + err.Error())
	}
	err = CreateHorsDoeuvres(c, &records)
	if err != nil {
		return errors.New("There was a problem creating the hors doeuvres records: " + err.Error())
	}
	return nil
}

func getIsTestEnv() bool {
	test_env_str, _ := os.LookupEnv("TEST_ENV")
	isTestEnv, _ := strconv.ParseBool(test_env_str)
	return isTestEnv
}

func getIsMockEnv() bool {
	test_env_str, _ := os.LookupEnv("USE_MOCK_DB")
	isTestEnv, _ := strconv.ParseBool(test_env_str)
	return isTestEnv
}

func checkTestEnv() error {
	isTestEnv := getIsTestEnv()
	if !isTestEnv {
		return errors.New("TEST_ENV was either not found, not defined or set to 'false' - must be set to 'true' for test database operations")
	}
	return nil
}

// Seeds test_db with test data defined in the /test-fixtures directory
func SeedTestData() {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	err := loadTestEntrees(ctx)
	if err != nil {
		log.Println(err.Error())
	}
	err = loadTestHorsDoeuvres(ctx)
	if err != nil {
		log.Println(err.Error())
	}
	err = checkTestEnv()
	if err != nil {
		log.Println(err.Error())
	}
	err = loadTestUsers(ctx)
	if err != nil {
		log.Println(err.Error())
	}
}

// Drops test_db
//
// test_db does not exist before or after tests. As a result, the "production"
// client controls the setup and teardown of the test DB. While the tests
// are running, the gorm client is configured to run commands on the test_db
// database.
func CreateTestDB() error {
	checkTestEnv()
	result := db.Exec("CREATE DATABASE test_db")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Drops test_db
//
// The gorm db connection MUST be connected to a database other than test_db when this is
// run, otherwise the command will fail since the client is connected to the DB it's trying
// to drop.
func DropTestDB() error {
	checkTestEnv()
	result := db.Exec("DROP DATABASE test_db")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Helper method to switch the underlying DB that the Gorm client is connected to.
//
// When switching the connection, DB().Close() MUST be called on the previous connection. This
// matters most when resetting test_db since it's dropped as part of the reset process. If the
// connection is not closed before attempting to drop the DB, the operation will fail and indicate
// there is still an active connection.
func SwitchConnectedDB(dbName string) error {
	// Close() is not normally required; however, we need to close the prior connection
	// so there is no longer a live connection to the test_db (otherwise, we can't DROP
	// it).
	conn, err := db.DB()
	if err != nil {
		return err
	}

	conn.Close()

	// Re-establish connection to DB using "production" DB name so we can drop the test DB
	dbConnectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("PGSQL_HOST"),
		os.Getenv("PGSQL_USER"),
		os.Getenv("PGSQL_PASSWORD"),
		dbName,
		os.Getenv("PGSQL_PORT"),
		os.Getenv("PGSQL_TIMEZONE"))

	db, err = gorm.Open(postgres.Open(dbConnectionString), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Panic("There was a problem connecting to the database: ", err.Error())
		return err
	}
	return nil
}

// Convenvience method to reset the test DB
func ResetAndConnectToTestDb() {
	SwitchConnectedDB(os.Getenv("PGSQL_DBNAME"))
	DropTestDB()
	CreateTestDB()
	SwitchConnectedDB("test_db")
	err := Migrate()
	if err != nil {
		log.Panic("There was a problem migrating the schema: ", err.Error())
	}
}
