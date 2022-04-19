package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/GrandOichii/mtgsdk"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TokenPath = "TOKEN"

	MaxCards = 20
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func readToken() (string, error) {
	b, err := os.ReadFile(TokenPath)
	return string(b), err
}

func main() {
	token, err := readToken()
	checkErr(err)
	bot, err := tgbotapi.NewBotAPI(token)
	checkErr(err)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	fmt.Println("Bot started")

	for update := range updates {
		if update.InlineQuery != nil && update.InlineQuery.Query != "" {
			q := update.InlineQuery.Query
			rm := []interface{}{}
			cards, err := mtgsdk.GetCards(map[string]string{mtgsdk.CardNameKey: q})
			checkErr(err)
			if len(cards) > MaxCards {
				cards = cards[:MaxCards]
			}
			for i, card := range cards {
				if card.ImageUris.Small == "" || card.ImageUris.Large == "" {
					continue
				}

				// r := tgbotapi.NewInlineQueryResultPhoto(strconv.Itoa(i), card.ImageUris.Large)
				// r := tgbotapi.NewInlineQueryResultArticle(strconv.Itoa(i), card.Name, "")
				html := fmt.Sprintf("<a href=\"%s\">%s</a>", card.ImageUris.Large, card.Name)
				r := tgbotapi.NewInlineQueryResultArticleHTML(strconv.Itoa(i), card.Name, html)
				r.ThumbURL = card.ImageUris.Small
				rm = append(rm, r)
			}
			ic := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				CacheTime:     0,
				Results:       rm,
			}

			_, err = bot.Request(ic)
			checkErr(err)
			continue
		}
	}
}
