package db

import "log"

func CheckBan(chatID int64) int {
	stmt, err := db.Prepare("SELECT is_banned FROM bot_users WHERE chat_id = ?")
	defer stmt.Close()
	if err != nil {
		log.Println("Ошибка: не удалось сформировать запрос на проверку бана.")
	}

	rows, err := stmt.Query(chatID)
	defer rows.Close()
	if err != nil {
		log.Println("Ошибка: не удалось выполнить SELECT забаненного юзера.")
	}

	for rows.Next() {
		var isBanned int

		err = rows.Scan(&isBanned)
		if err != nil {
			log.Println("Ошибка: при использовании rows.Scan произошла ошибка.")
		} else if isBanned == 1 {
			return 1
		}
	}

	return 0
}

func AddBan(chatID int64) error {
	stmt, err := db.Prepare("UPDATE bot_users SET is_banned = 1 WHERE chat_id = ?")
	defer stmt.Close()
	if err != nil {
		log.Println("Ошибка: не удалось сформировать запрос на добавление бана.")
		return err
	}

	_, err = stmt.Query(chatID)
	if err != nil {
		log.Println("Ошибка: не удалось выполнить запрос на бан юзера.")
		return err
	}

	return nil
}

func RemoveBan(chatID int64) error {
	stmt, err := db.Prepare("UPDATE bot_users SET is_banned = 0 WHERE chat_id = ?")
	defer stmt.Close()
	if err != nil {
		log.Println("Ошибка: не удалось сформировать запрос на удаление бана.")
		return err
	}

	_, err = stmt.Query(chatID)
	if err != nil {
		log.Println("Ошибка: не удалось выполнить запрос на разбан юзера.")
		return err
	}

	return nil
}
