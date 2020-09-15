package db

import (
	"fmt"
	"log"
)

type Chat struct {
	Id     int
	ChatId int64
}

func (chat Chat) String() string {
	return fmt.Sprintf("Id:%d chatid:%d", chat.Id, chat.ChatId)
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
		err = rows.Scan(&id, &chatid)
		if err != nil {
			log.Printf("error while scanning chat table: %s", err)
		}
		result = append(result, Chat{
			Id:     id,
			ChatId: chatid,
		})
	}
	return result
}

func InsertChat(chat Chat) {
	insertSQL := "insert into chat(chatid) values(?)"
	statement, err := Database.Prepare(insertSQL)
	if err != nil {
		log.Fatalf("failed to insert into chat, error: %s", err)
	}
	_, err = statement.Exec(chat.ChatId)
	if err != nil {
		log.Printf("failed to insert chatid %d, error: %s", chat.ChatId, err)
	}
}

func DeleteChat(chatid int64) {
	insertSQL := "delete from chat where chatid=?"
	statement, err := Database.Prepare(insertSQL)
	if err != nil {
		log.Fatalf("failed to delete chat with chatid %d, error: %s", chatid, err)
	}
	_, err = statement.Exec(chatid)
	if err != nil {
		log.Printf("failed to insert chatid %d, error: %s", chatid, err)
	}
}
