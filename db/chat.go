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
		log.Fatalf("failed to query table chat, error: %s", err)
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
