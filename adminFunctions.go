package main

import (
	"DevByBot-2.0/db"
)

func BanUser(chatID int64) error {
	if err := db.AddBan(chatID); err != nil {
		return err
	}

	return nil
}

func UnbanUser(chatID int64) error {
	if err := db.RemoveBan(chatID); err != nil {
		return err
	}

	return nil
}

func AlertUsers(msg string) {
	chatIDs := db.SelectAllChatID()
	for index := 0; index < chatIDs.Size(); index++ {
		chatID, _ := chatIDs.Get(index)
		sendSimpleMsg(int64(chatID.(int)), msg)
	}
}
