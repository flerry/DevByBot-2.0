package parsers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mmcdole/gofeed"
	"log"
)

func ParseEvents(chatID int64) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://events.dev.by/rss")

	if err != nil {
		log.Println("Не удалось спарсить события.")
		msg := tgbotapi.NewMessage(chatID, "Извините, мне не удалось получить список ближайших событий. Скорее всего это - проблема сайта.")
		bot.Send(msg)
		return
	}

	titles := make([]string, 7, 7)
	links := make([]string, 7, 7)

	msg := tgbotapi.NewMessage(chatID, "Ближайшие события:")
	for element := 0; element <= 6; element++ {
		titles[element] = feed.Items[element].Title
		links[element] = feed.Items[element].Link
	}

	msg.ReplyMarkup = setInlineEvents(links, titles)
	bot.Send(msg)
}

func setInlineEvents(links []string, titles []string) tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("#1 "+titles[0], links[0])),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("#2 "+titles[1], links[1])),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("#3 "+titles[2], links[2])),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("#4 "+titles[3], links[3])),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("#5 "+titles[4], links[4])),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("#6 "+titles[5], links[5])),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("#7 "+titles[6], links[6])))

	return kb
}
