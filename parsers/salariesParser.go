package parsers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func ParseSalaries(chatID int64) {
	doc, err := goquery.NewDocument("https://salaries.dev.by")
	if err != nil {
		log.Println("Не удалось спарсить зарплаты.")
		msg := tgbotapi.NewMessage(chatID, "Извините, мне не удалось получить список зарплат. Скорее всего это - проблема сайта.")
		bot.Send(msg)
		return
	}

	var average, median string
	doc.Find("#dev-salaries > div.block.data-summary-info > div:nth-child(2) > strong").Each(func(i int, s *goquery.Selection) {
		average = s.Text()
	})

	doc.Find("#dev-salaries > div.block.data-summary-info > div:nth-child(3) > strong").Each(func(i int, s *goquery.Selection) {
		median = s.Text()
	})
	msg := tgbotapi.NewMessage(chatID, "Зарплата в IT:")

	msg.ReplyMarkup = setInlineSalaries(average, median)
	bot.Send(msg)
}

func setInlineSalaries(average string, median string) tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Средняя: "+average+"$", "1"),
			tgbotapi.NewInlineKeyboardButtonData("Медиана: "+median+"$", "1")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Узнать подробнее", "https://salaries.dev.by")))

	return kb
}
