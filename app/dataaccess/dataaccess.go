package dataaccess

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

const (
	SOURCE_SLACK = "slack"
	SOURCE_WEB   = "web"
)

func CreateAck(db *sql.DB, message string, senderEmail string, source string) error {
	//TODO handle case where msg is empty
	//TODO trim spaces in message
	//TODO handle upserting slack info
	query := "INSERT INTO users (email, updated_at) values ($1, current_timestamp) ON CONFLICT (email) DO UPDATE SET updated_at = current_timestamp RETURNING id"
	rows, err := db.Query(query, senderEmail)
	//defer rows.Close()
	if err != nil {
		log.Println(err) //FIXME handle this better
		return err
	}

	var userId string
	rows.Next()
	err = rows.Scan(&userId)
	if err != nil {
		log.Println(err)
	}

	query = "INSERT INTO acks (msg, user_id, source, updated_at) values ($1, $2, $3, current_timestamp)"
	_, err = db.Exec(query, message, userId, source)
	if err != nil {
		log.Println(err) //FIXME handle this better
		return err
	}
	return err
}

// empty senderEmail string queries for all acks
func FetchAcks(db *sql.DB, senderEmail string) []string {
	log.Println("called FetchAcks()")
	messages := make([]string, 0)
	var rows *sql.Rows
	var err error
	var query string

	//FIXME figure out how to properly pass timestamps to db.Query instead of concatenating a string here
	curDateString := fmt.Sprintf(time.Now().Format("2006-01-02"))
	if senderEmail == "" {
		log.Println("fetching all acks, last 7 days")
		query = "select msg from acks where created_at > ('" + curDateString + "' - 7) order by updated_at desc"
		rows, err = db.Query(query)
	} else {
		log.Println("fetching acks by user email: " + senderEmail)
		query = "select msg from acks a, users u where u.email = $1 and a.user_id = u.id and a.created_at > ('" + curDateString + "' - 7) order by a.updated_at desc"
		rows, err = db.Query(query, senderEmail)
	}

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var message string
	log.Println("fetching acks:")
	for rows.Next() {
		err = rows.Scan(&message)
		if err != nil {
			log.Println(err)
		}

		log.Println(message)
		messages = append(messages, message)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}

	return messages
}
