package main

import (
	"DevByBot-2.0/db"
	"DevByBot-2.0/parsers"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var bot *tgbotapi.BotAPI

func init() {
	formatDate := time.Now().Local().Format("2006-01-02")
	formatTime := time.Now().Local().Format("15:00:00")

	f, err := os.Create("log" + formatDate + "|" + formatTime + ".txt")

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	f.Write([]byte(formatTime))

	log.SetOutput(f)

	db.Setup()

	bot, _ = tgbotapi.NewBotAPI("314564211:AAH59sKgMcht-F_sVevp3jGXLo9j2VELRqg")
	if len(bot.Token) == 0 {
		log.Fatal("Ошибка: не удалось инициализировать бота.")
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	parsers.SetBot(bot)

	go parsers.ParseFeed("default")
}

func main() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Fatal("Ошибка: не удалось получить обновления.")
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[id %d][%s] %s", update.Message.Chat.ID, update.Message.From.UserName, update.Message.Text)

		if db.CheckBan(update.Message.Chat.ID) == 0 {
			switch {
			case update.Message.Text == "/start":
				{
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,
						"Привет, "+update.Message.Chat.UserName+"! Я - бот ресурса Dev.by. "+
							"Получайте ИТ-новости одновременно с их выходом на сайте и пользуйтесь дополнительными возможностями с помощью кнопок.")
					msg.ReplyMarkup = setSubscribeKeyboard()

					bot.Send(msg)
				}
			case update.Message.Text == "🔔 Subscribe":
				{
					go db.InsertChatID(update.Message.Chat.ID)

					var infoText string

					if update.Message.Chat.UserName == "" {
						infoText = "Ваша подписка успешно оформлена! " +
							"Чтобы удалить подписку, нажмите на \"Unsubscribe\""
					} else {
						infoText = update.Message.Chat.UserName + ", Ваша подписка успешно оформлена! " +
							"Чтобы удалить подписку, нажмите на \"Unsubscribe\""
					}

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, infoText)
					msg.ReplyMarkup = setUnsubscribeKeyboard()
					bot.Send(msg)
				}
			case update.Message.Text == "🔕 Unsubscribe":
				{
					go db.RemoveChatID(update.Message.Chat.ID)

					var infoText string

					if update.Message.Chat.UserName == "" {
						infoText = "Ваша подписка успешно удалена! " +
							"Чтобы оформить подписку, нажмите на \"Subscribe\""
					} else {
						infoText = update.Message.Chat.UserName + ", Ваша подписка успешно удалена! " +
							"Чтобы оформить подписку, нажмите на \"Subscribe\""
					}

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, infoText)
					msg.ReplyMarkup = setSubscribeKeyboard()
					bot.Send(msg)
				}
			case update.Message.Text == "✨ Events":
				{
					go parsers.ParseEvents(update.Message.Chat.ID)
				}
			case update.Message.Text == "💳 Salaries":
				{
					go parsers.ParseSalaries(update.Message.Chat.ID)
				}
			case update.Message.Text == "🤝 Community":
				{
					go sendSocialMsg(update.Message.Chat.ID)
				}
			case strings.Contains(update.Message.Text, "alert4532"):
				{
					alertMsgText := strings.Replace(update.Message.Text, "alert4532 ", "", -1)
					AlertUsers(alertMsgText)
					sendSimpleMsg(update.Message.Chat.ID,
						fmt.Sprintf("%s: [%s] %s", "Сообщение", alertMsgText, "отправлено."))
				}

			case strings.Contains(update.Message.Text, "unban4532"):
				{
					unBannedChatID, _ := strconv.ParseInt(strings.Replace(update.Message.Text, "unban4532 ", "", -1),
						10, 64)
					if err := UnbanUser(unBannedChatID); err != nil {
						sendSimpleMsg(update.Message.Chat.ID, "Извините, но пользователя не удалось разабанить.")
					} else {
						sendSimpleMsg(update.Message.Chat.ID,
							fmt.Sprintf("%s: [id%d] %s", "Пользователь", unBannedChatID, "разбанен."))
					}
				}
			case strings.Contains(update.Message.Text, "ban4532"):
				{
					bannedChatID, _ := strconv.ParseInt(strings.Replace(update.Message.Text, "ban4532 ", "", -1),
						10, 64)
					if err := BanUser(bannedChatID); err != nil {
						sendSimpleMsg(update.Message.Chat.ID, "Извините, но пользователя не удалось забанить.")
					} else {
						sendSimpleMsg(update.Message.Chat.ID,
							fmt.Sprintf("%s: [id%d] %s", "Пользователь", bannedChatID, "забанен."))
					}
				}
			default:
				sendSimpleMsg(update.Message.Chat.ID, "Извините, но я не знаю, что на это ответить.")
			}
		} else {
			sendSimpleMsg(update.Message.Chat.ID, "Извините, но Вы временно заблокированы.")
		}
	}
}

func sendSimpleMsg(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "markdown"
	msg.DisableWebPagePreview = true
	bot.Send(msg)
}

func setSubscribeKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🔔 Subscribe"),
	),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✨ Events"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("💳 Salaries"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🤝 Community"),
		),
	)
}

func setUnsubscribeKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🔕 Unsubscribe"),
	),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✨ Events"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("💳 Salaries"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🤝 Community"),
		),
	)
}

func sendSocialMsg(chatID int64) {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("SLACK", "https://devby.slack.com")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("VK", "https://vk.com/devby")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("FACEBOOK", "https://www.facebook.com/devbyby")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("TWITTER", "https://twitter.com/devby")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("INSTAGRAM", "https://twitter.com/devby")))

	msg := tgbotapi.NewMessage(chatID, "Cообщества dev.by:")
	msg.DisableWebPagePreview = true
	msg.ReplyMarkup = kb

	bot.Send(msg)
}
