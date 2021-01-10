package db

import (
	"fmt"
	"log"
)

const UserNameRoleAdmin = "admin"
const UserNameRoleReader = "reader"

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
		var name, role string
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

/**
  Insert a row into the username table. Returns the lastInsertId of the insert operation.
*/
func InsertUserName(userName UserName) int64 {
	insertSQL := "insert into username(name,role) values(?,?)"
	statement, err := Database.Prepare(insertSQL)
	defer statement.Close()
	if err != nil {
		log.Printf("failed to prepare stmt for insert into username, error: %s", err)
		return 0
	} else {
		result, err := statement.Exec(userName.Name, userName.Role)
		if err != nil {
			log.Printf("failed to insert username %d, error: %s", userName.Name, err)
			return 0
		} else {
			lastInsertId, err := result.LastInsertId()
			if err == nil {
				return lastInsertId
			} else {
				log.Printf("no username row was inserted, err: %s", err)
				return 0
			}
		}
	}
}

func DeleteUserName(name string) bool {
	deleteSQL := "delete from username where name=?"
	statement, err := Database.Prepare(deleteSQL)
	defer statement.Close()
	if err != nil {
		log.Printf("failed to prepare stmt for delete username with name %s, error: %s", name, err)
		return false
	} else {
		result, err := statement.Exec(name)
		if err != nil {
			log.Printf("failed to delete username %s, error: %s", name, err)
			return false
		} else {
			rowsAffected, _ := result.RowsAffected()
			return rowsAffected == 1
		}
	}
}
