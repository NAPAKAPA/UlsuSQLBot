package main

import (
	"github.com/education-bot/pkg/core"
	"github.com/education-bot/pkg/core/models"
	"github.com/education-bot/pkg/core/routes"
	"github.com/education-bot/pkg/utils"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var (
	userForm = tgbot.NewInlineKeyboardMarkup(
		tgbot.NewInlineKeyboardRow(
			tgbot.NewInlineKeyboardButtonData("Регистрация", routes.UserRegistrationCmd),
			tgbot.NewInlineKeyboardButtonData("Обновить данные", routes.UserUpdateCmd),
		),
	)
)

func main() {
	bot, err := tgbot.NewBotAPI(utils.GetEnv("TG_BOT_TOKEN", "6794171875:AAH1Ti1ZXf4R9trfw6D-yruc691hfFRXyuU"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	update := tgbot.NewUpdate(0)
	update.Timeout = 60
	updatesChan := bot.GetUpdatesChan(update)
	for update := range updatesChan {
		err, cont, user := routes.Authentication(bot, update)
		if err != nil && !cont {
			msg := tgbot.NewMessage(update.Message.Chat.ID, "Вам нужно пройти регистрацию")
			msg.ReplyMarkup = userForm
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			continue
		}
		if user == nil {
			var userId int64
			if update.Message != nil {
				userId = update.Message.From.ID
			}
			if update.CallbackQuery != nil {
				userId = update.CallbackQuery.From.ID
			}
			userR := core.User{Id: userId}
			uReq := models.UserRequest{}
			user, err = uReq.Get(userR.Key())
		}
		/*если пользователь не авторизован - пропускать следующий шаг*/
		if user != nil && update.CallbackQuery != nil {
			err, cont = routes.SqlTasks(bot, update)
			continue
		}
		if user != nil {
			err, cont = routes.SqlTasks(bot, update)
			if err != nil {
				continue
			}
		}
	}
}
