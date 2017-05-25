package db

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"log"
)

func SelectAllChatID() (listID arraylist.List) {
	rows, err := db.Query("SELECT chat_id FROM bot_users")
	defer rows.Close()
	if err != nil {
		log.Fatal("Ошибка: не удалось получить список chat_id.")
	}

	chatIDs := arraylist.New()
	for rows.Next() {
		var chat_id int

		err = rows.Scan(&chat_id)
		if err != nil {
			log.Println("Ошибка: при использовании rows.Scan произошла ошибка.")
		} else {
			chatIDs.Add(chat_id)
		}
	}

	return *chatIDs
}

func InsertChatID(chatID int64) error {
	stmt, err := db.Prepare("INSERT INTO bot_users (chat_id) VALUES (?)")
	defer stmt.Close()
	if err != nil {
		log.Println("Ошибка: не удалось сформировать запрос на добавление chat_id.")
		return err
	}

	_, err = stmt.Query(chatID)
	if err != nil {
		log.Println("Ошибка: не удалось выполнить добавление chat_id или chat_id уже есть в бд.")
		return err
	}

	return nil
}

func RemoveChatID(chatID int64) error {
	stmt, err := db.Prepare("DELETE FROM bot_users WHERE chat_id = ?")
	defer stmt.Close()
	if err != nil {
		log.Println("Ошибка: не удалось сформировать запрос на удаление chat_id.")
		return err
	}

	_, err = stmt.Query(chatID)
	if err != nil {
		log.Println("Ошибка: не удалось выполнить удаление chat_id.")
		return err
	}

	return nil
}
