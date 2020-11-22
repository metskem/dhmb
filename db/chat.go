package db

import (
	"fmt"
	"log"
)

type Chat struct {
	Id     int
	ChatId int64
	Name   string
}

func (chat Chat) String() string {
	return fmt.Sprintf("Id:%d chatid:%d Name:%s", chat.Id, chat.ChatId, chat.Name)
}

func GetChats() []Chat {
	rows, err := Database.Query("select * from chat order by chatid", nil)
	if err != nil {
		log.Printf("failed to query table chat, error: %s", err)
	}
	defer rows.Close()
	var result []Chat
	for rows.Next() {
		var id int
		var chatid int64
		var name string
		err = rows.Scan(&id, &chatid, &name)
		if err != nil {
			log.Printf("error while scanning chat table: %s", err)
		}
		result = append(result, Chat{
			Id:     id,
			ChatId: chatid,
			Name:   name,
		})
	}
	return result
}

/**
  Insert a row into the chat table using the given chat. Returns the lastInsertId of the insert operation.
*/
func InsertChat(chat Chat) int64 {
	insertSQL := "insert into chat(chatid, name) values(?,?)"
	statement, err := Database.Prepare(insertSQL)
	defer statement.Close()
	if err != nil {
		log.Printf("failed to prepare stmt for insert into chat, error: %s", err)
		return 0
	} else {
		result, err := statement.Exec(chat.ChatId, chat.Name)
		if err != nil {
			log.Printf("failed to insert chatid %d, name %s, error: %s", chat.ChatId, chat.Name, err)
			return 0
		} else {
			lastInsertId, err := result.LastInsertId()
			if err == nil {
				return lastInsertId
			} else {
				log.Printf("no chat row was inserted, err: %s", err)
				return 0
			}
		}
	}
}

func DeleteChat(chatid int64) bool {
	deleteSQL := "delete from chat where chatid=?"
	statement, err := Database.Prepare(deleteSQL)
	defer statement.Close()
	if err != nil {
		log.Printf("failed to prepare stmt for delete chat with chatid %d, error: %s", chatid, err)
		return false
	} else {
		result, err := statement.Exec(chatid)
		if err != nil {
			log.Printf("failed to delete chatid %d, error: %s", chatid, err)
			return false
		} else {
			rowsAffected, _ := result.RowsAffected()
			return rowsAffected == 1
		}
	}
}
