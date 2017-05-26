package parsers

import (
	"DevByBot-2.0/db"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mmcdole/gofeed"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Post struct {
	Title     string
	Link      string
	ImageLink string
}

var bot *tgbotapi.BotAPI

func SetBot(newBot *tgbotapi.BotAPI) {
	bot = newBot
}

func ParseFeed(nowLink string) {
	oldLink := nowLink

	fp := gofeed.NewParser()
	for {
		feed, err := fp.ParseURL("https://dev.by/rss")

		if err != nil {
			log.Println("Не удалось спарсить rss.")
			ParseFeed(oldLink)
		}

		link := feed.Items[0].Link

		if link != oldLink {
			post := Post{
				Title:     feed.Items[0].Title,
				Link:      link,
				ImageLink: feed.Items[0].Enclosures[0].URL,
			}

			listID := db.SelectAllChatID()
			SendNews(&post, listID)
			oldLink = link
		}
		time.Sleep(time.Minute * 1)
	}
}

func SendNews(post *Post, chatIDs arraylist.List) {
	response, e := http.Get(post.ImageLink)
	if e != nil {
		log.Fatal(e)
	}

	defer response.Body.Close()

	var name string
	if strings.Contains(post.ImageLink, "jpg") {
		name = "post_photo.jpg"
	} else {
		name = "post_photo.png"
	}

	file, err := os.Create(name)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	for index := 0; index < chatIDs.Size(); index++ {
		chatID, _ := chatIDs.Get(index)
		caption := post.Title

		msg := tgbotapi.NewPhotoUpload(int64(chatID.(int)), name)
		msg.Caption = caption

		msg.ReplyMarkup = setInlineURL(post.Link, "Читать далее")

		bot.Send(msg)
	}
}

func setInlineURL(URL string, inlineText string) tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(inlineText, URL)))

	return kb
}
