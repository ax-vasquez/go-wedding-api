package models

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
)

func loadTestUsers() error {
	users := []User{}
	userInvitees := []User{}
	userDataFile, err := os.ReadFile("../test-fixtures/users.json")
	if err != nil {
		return errors.New("There was a problem loading test user data from ./test-fixtures/users.json: " + err.Error())
	}
	userInviteeDataFile, err := os.ReadFile("./../test-fixtures/invitees.json")
	if err != nil {
		return errors.New("There was a problem loading test user data from ./test-fixtures/invitees.json: " + err.Error())
	}
	err = json.Unmarshal(userDataFile, &users)
	if err != nil {
		return errors.New("There was a problem unmarshaling the JSON from file ./test-fixtures/users.json: " + err.Error())
	}
	err = json.Unmarshal(userInviteeDataFile, &userInvitees)
	if err != nil {
		return errors.New("There was a problem unmarshaling the JSON from file ./test-fixtures/invitees.json: " + err.Error())
	}
	_, err = CreateUsers(&users)
	if err != nil {
		return errors.New("There was a problem creating test user records: " + err.Error())
	}
	_, err = CreateUsers(&userInvitees)
	if err != nil {
		return errors.New("There was a problem creating test user invitee records: " + err.Error())
	}
	return nil
}

func loadTestUserInviteeRelationships() error {
	records := []UserUserInvitee{}
	recordsFile, err := os.ReadFile("../test-fixtures/user_user_invitees.json")
	if err != nil {
		return errors.New("There was a problem loading test user data from ./test-fixtures/user_user_invitees.json: " + err.Error())
	}
	err = json.Unmarshal(recordsFile, &records)
	if err != nil {
		return errors.New("There was a problem unmarshaling the JSON from file ./test-fixtures/user_user_invitees.json: " + err.Error())
	}
	err = CreateUserUserInvitees(records)
	if err != nil {
		return errors.New("There was a problem creating the test user user invitee records: " + err.Error())
	}
	return nil
}

func loadTestEntrees() error {
	records := []Entree{}
	recordsFile, err := os.ReadFile("../test-fixtures/entrees.json")
	if err != nil {
		return errors.New("There was a problem loading test user data from ./test-fixtures/entrees.json: " + err.Error())
	}
	err = json.Unmarshal(recordsFile, &records)
	if err != nil {
		return errors.New("There was a problem unmarshaling the JSON from file ./test-fixtures/entrees.json: " + err.Error())
	}
	_, err = CreateEntrees(&records)
	if err != nil {
		return errors.New("There was a problem creating the test entree records: " + err.Error())
	}
	return nil
}

func loadTestHorsDoeuvres() error {
	records := []HorsDoeuvres{}
	recordsFile, err := os.ReadFile("../test-fixtures/hors_doeuvres.json")
	if err != nil {
		return errors.New("There was a problem loading test user data from ./test-fixtures/hors_doeuvres.json: " + err.Error())
	}
	err = json.Unmarshal(recordsFile, &records)
	if err != nil {
		return errors.New("There was a problem unmarshaling the JSON from file ./test-fixtures/hors_doeuvres.json: " + err.Error())
	}
	_, err = CreateHorsDoeuvres(&records)
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

func checkTestEnv() error {
	isTestEnv := getIsTestEnv()
	if !isTestEnv {
		return errors.New("TEST_ENV was either not found, not defined or set to 'false' - must be set to 'true' for test database operations")
	}
	return nil
}

// Seeds test_db with test data defined in the /test-fixtures directory
func SeedTestData() {
	err := loadTestEntrees()
	if err != nil {
		log.Println(err.Error())
	}
	err = loadTestHorsDoeuvres()
	if err != nil {
		log.Println(err.Error())
	}
	err = checkTestEnv()
	if err != nil {
		log.Println(err.Error())
	}
	err = loadTestUsers()
	if err != nil {
		log.Println(err.Error())
	}
	err = loadTestUserInviteeRelationships()
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
