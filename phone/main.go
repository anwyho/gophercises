package main

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

type psqlInfo struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

func (p *psqlInfo) String() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.host, p.port, p.user, p.password, p.dbname, "disable")
}

const (
	host     = "localhost"
	port     = 5432
	user     = "anthonyho"
	password = "FNN9nA2MEKuHDK3q"
	dbname   = "gophercises_phone"
)

var digitsRe *regexp.Regexp

func init() {
	digitsRe = regexp.MustCompile(`\D`)
}

func normalize(d string) string {
	return digitsRe.ReplaceAllString(d, "")
}

func createDB(db *sql.DB, name string) (err error) {
	_, err = db.Exec("CREATE DATABASE " + name)
	return
}

func resetDB(db *sql.DB, name string) (err error) {
	_, err = db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createPhoneTable(db *sql.DB) (err error) {
	statement := `
        CREATE TABLE IF NOT EXISTS phone_numbers (
            id SERIAL,
            value VARCHAR(64)
        )`
	_, err = db.Exec(statement)
	return err
}

func insertPhoneNumber(db *sql.DB, number string) (id int, err error) {
	s := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	if err = db.QueryRow(s, number).Scan(&id); err != nil {
		id = -1
	}
	return
}

func getPhoneNumber(db *sql.DB, id int) (number string, err error) {
	s := "SELECT * FROM phone_numbers WHERE id=$1"
	if err = db.QueryRow(s, id).Scan(&id, &number); err != nil {
		number = ""
	}
	return
}

func findPhoneNumber(db *sql.DB, number string) (*phoneNumber, error) {
	var p phoneNumber
	s := "SELECT * FROM phone_numbers WHERE value=$1"
	if err := db.QueryRow(s, number).Scan(&p.id, &p.number); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

type phoneNumber struct {
	id     int
	number string
}

func updatePhoneNumber(db *sql.DB, pn phoneNumber) (err error) {
	s := "UPDATE phone_numbers SET value=$2 WHERE id=$1"
	_, err = db.Exec(s, pn.id, pn.number)
	return
}

func deletePhoneNumber(db *sql.DB, id int) (err error) {
	s := "DELETE FROM phone_numbers WHERE id=$1"
	_, err = db.Exec(s, id)
	return
}

func getAllPhoneNumbers(db *sql.DB) (phoneNumbers []phoneNumber, err error) {
	rows, err := db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p phoneNumber
		if err = rows.Scan(&p.id, &p.number); err != nil {
			phoneNumbers = nil
			return
		}
		phoneNumbers = append(phoneNumbers, p)
	}
	if err = rows.Err(); err != nil {
		phoneNumbers = nil
	}
	return
}

func main() {
	var err error
	var db *sql.DB
	pi := psqlInfo{host, port, user, password, "postgres"}

	db, err = sql.Open("postgres", pi.String())
	if err != nil {
		panic(err)
	}
	if err = resetDB(db, dbname); err != nil {
		panic(err)
	}
	db.Close()

	pi.dbname = dbname
	db, err = sql.Open("postgres", pi.String())
	if err != nil {
		panic(err)
	}
	if err = createPhoneTable(db); err != nil {
		panic(err)
	}
	addNumbers := []string{
		"+1(234)567-8901",
		"11111111111",
		"+1(234567-8901",
		"+1(234)567-8901",
		"+1(111)111-1111",
		"11111111111",
	}
	for _, n := range addNumbers {
		if _, err = insertPhoneNumber(db, n); err != nil {
			panic(err)
		}
	}
	var phoneNumbers []phoneNumber
	if phoneNumbers, err = getAllPhoneNumbers(db); err != nil {
		panic(err)
	}
	for _, pn := range phoneNumbers {
		fmt.Printf("Working on... %+v\n", pn)
		number := normalize(pn.number)
		if existingPn, err := findPhoneNumber(db, number); err != nil {
			panic(err)
		} else {
			if existingPn != nil && existingPn.id != pn.id {
				fmt.Println("deleting number ", pn.number)
				err = deletePhoneNumber(db, pn.id)
			} else {
				fmt.Println("updating number ", pn.number)
				pn.number = number
				err = updatePhoneNumber(db, pn)
			}
			if err != nil {
				panic(err)
			}
		}
	}
	if phoneNumbers, err = getAllPhoneNumbers(db); err != nil {
		panic(err)
	}
	fmt.Println("normalized numbers: ")
	for _, pn := range phoneNumbers {
		fmt.Println(pn.number)
	}
}

// func normalize(d string) string {
// 	var buf bytes.Buffer
// 	for _, c := range d {
// 		if c >= '0' && c <= '9' {
// 			buf.WriteRune(c)
// 		}
// 	}
// 	return buf.String()
// }
