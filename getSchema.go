package post05

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

var (
	Hostname = ""
	Port     = 5432
	Username = ""
	Password = ""
	Database = ""
)

type UserData struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

func openConnection() (*sql.DB, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Hostname, Port, Username, Password, Database)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
func exists(username string) int {
	username = strings.ToLower(username)
	userid := -1
	db, err := openConnection()
	if err != nil {
		print(err)
		return userid
	}
	defer db.Close()
	query := fmt.Sprintf(`SELECT "id" FROM "users" WHERE username=%s`, username)
	rows, err := db.Query(query)
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			print("Scan Error", err)
			return userid
		}
		userid = id
	}
	defer rows.Close()
	return userid
}
func AddUser(d UserData) int {
	d.Username = strings.ToLower(d.Username)
	db, err := openConnection()
	if err != nil {
		print("ERROR while adding in openning", err)
		return -1
	}
	defer db.Close()
	userid := exists(d.Username)
	if userid != -1 {
		print("User already existed...")
		return -1
	}
	insertStatement := `INSERT INTO "users" ("username") VALUES ($1)`
	_, err = db.Exec(insertStatement, d.Username)
	if err != nil {
		print(err)
		return -1
	}
	userid = exists(d.Username)
	if userid == -1 {
		return -1
	}
	insertStatement = `INSERT INTO "userdata" ("userid","name","surname","description") VALUES ($1,$2,$3,$4)`
	_, err = db.Exec(insertStatement, userid, d.Username, d.Surname, d.Description)
	if err != nil {
		print("inserting into userdata table", err)
		return -1
	}
	return userid
}
func DeleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	statment := fmt.Sprintf(`SELECT "username" FROM "users" WHERE id=%d`, id)
	rows, err := db.Query(statment)
	if err != nil {
		return err
	}
	defer rows.Close()
	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return err
		}
	}
	if exists(username) != id {
		return fmt.Errorf("User with id=%d doesn't exist", id)
	}
	deleteStatement := `DELETE FROM "users" where id=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}
	deleteStatement = `DELETE FROM "userdata" WHERE userid=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}
	return nil
}
func ListUsers() ([]UserData, error) {
	Data := []UserData{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close()
	query := `SELECT "id","username","name","surname","description" 
			FROM "users","userdata"
			WHERE users.id=userdata.id`
	rows, err := db.Query(query)
	if err != nil {
		return Data, err
	}
	for rows.Next() {
		var id int
		var username string
		var name string
		var surname string
		var description string
		err = rows.Scan(&id, &username, &name, &surname, &description)
		if err != nil {
			return Data, err
		}
		temp := UserData{ID: id, Username: username, Name: name, Surname: surname, Description: description}
		Data = append(Data, temp)
	}
	defer rows.Close()
	return Data, nil
}
func UpdateUser(d UserData) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	userId := exists(d.Username)
	if userId == -1 {
		return errors.New("User doesn't exist")
	}
	d.ID = userId
	updateStatement := `UPDATE "userdata" set "name"=$1, "surname"=$2,"description"=$3
	 WHERE "userid"=$4`
	_, err = db.Exec(updateStatement, d.Name, d.Surname, d.Description, d.ID)
	if err != nil {
		return err
	}
	return nil
}

var print = fmt.Println

//finished
