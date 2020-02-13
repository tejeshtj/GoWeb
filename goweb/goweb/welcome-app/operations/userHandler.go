package operations

/*We import 4 important libraries
1. “net/http” to access the core go http functionality
2. “fmt” for formatting our text
3. “html/template” a library that allows us to interact with our html file.
4. "time" - a library for working with date and time.*/
import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func CreateTables() {
	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/golangweb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	{ // Create a new table
		query := `
		CREATE TABLE users (
			id INT AUTO_INCREMENT,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME,
			PRIMARY KEY (id)
		);`

		if _, err := db.Exec(query); err != nil {
			log.Fatal(err)
		}
	}
}

func InsertData(username, password string) int64 {
	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/golangweb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	existingId, role := QueryData(username, password)
	fmt.Println(role)
	createdAt := time.Now()

	if existingId > 0 {
		return 0
	} else {
		result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createdAt)
		id, err := result.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id)
		return id
	}

}

func QueryData(username, password string) (int, string) {

	var id int
	var role string
	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/golangweb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	query := "SELECT id,role FROM users WHERE username=? and password=?"
	if err := db.QueryRow(query, username, password).Scan(&id, &role); err != nil {
		//og.Fatal(err)
		return 0, ""
	}

	fmt.Println("from userhandler", id)
	return id, role
}

func AllData() {
	type user struct {
		id        int
		username  string
		password  string
		createdAt time.Time
	}
	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/golangweb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(`SELECT id, username, password, created_at FROM users`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var u user

		err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v", users)
}

func DeleteData() {
	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/golangweb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	_, error := db.Exec(`DELETE FROM users WHERE id = ?`, 1)
	if error != nil {
		log.Fatal(err)
	}
}
