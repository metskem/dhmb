package db

import (
	"fmt"
	"log"
)

type UserName struct {
	Id   int
	Name string
	Role string
}

func (userName UserName) String() string {
	return fmt.Sprintf("Id:%d Name:%s Role:%s", userName.Id, userName.Name, userName.Role)
}

func GetUserNames() []UserName {
	rows, err := Database.Query("select * from username order by name", nil)
	if err != nil {
		log.Printf("failed to query table username, error: %s", err)
	}
	defer rows.Close()
	var result []UserName
	for rows.Next() {
		var id int
		var name string
		var role string
		err = rows.Scan(&id, &name, &role)
		if err != nil {
			log.Printf("error while scanning username table: %s", err)
		}
		result = append(result, UserName{
			Id:   id,
			Name: name,
			Role: role,
		})
	}
	return result
}
