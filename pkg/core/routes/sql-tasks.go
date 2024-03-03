package routes

import (
	"errors"
	"fmt"
	"github.com/education-bot/pkg/database"
	"github.com/education-bot/pkg/utils"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"reflect"
)

const (
	UserTestingCreateTableCmd      = "user-test-create-table"
	UserTestingAddColumnToTableCmd = "user-test-add-column-to-table"

	UserExecuteCmd = "user-execute-cmd"
)

var (
	sqlDriver  = database.MsSql{}
	TestingCmd = map[string]interface{}{
		UserTestingCreateTableCmd:      nil,
		UserTestingAddColumnToTableCmd: nil,
	}
	TCmds = []string{
		UserTestingCreateTableCmd,
		UserTestingAddColumnToTableCmd}
)

func SqlTasks(bot *tg.BotAPI, update tg.Update) (err error, cont bool) {
	var userId int64
	if update.CallbackQuery != nil {
		userId = update.CallbackQuery.From.ID
		cmd := update.CallbackQuery.Data
		if _, ok := TestingCmd[cmd]; !ok {
			cont = true
			return
		}
		switch cmd {
		default:
			err = errors.New("cmd not found")
			return
		case UserTestingCreateTableCmd:
			msg := tg.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Напишите скрипт создания таблицы с названием [teachers] используя базу данных [db%d] и схему [dbo] и укажите обязательные поля:\n"+
				"Id INT\n"+
				"Name NVARCHAR(50)", userId))
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			cont = true
			errExec := sqlDriver.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS db%d;", userId))
			if errExec != nil {
				return errExec, false
			}
			errExec = sqlDriver.Exec(fmt.Sprintf("CREATE DATABASE db%d;", userId))
			if errExec != nil {
				return errExec, false
			}
			actions[userId] = UserExecuteCmd
			cont = true
			return
		case UserTestingAddColumnToTableCmd:
			msg := tg.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Напишите скрипт добавления колонки в таблицу с названием [homes] используя базу данных [db%d] и схему [dbo] колонка:\n [username], тип данных колонки [NVARCHAR(50)]", userId))
			if _, err := bot.Send(msg); err != nil {
				return err, false
			}
			errExec := sqlDriver.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS db%d;", userId))
			if errExec != nil {
				return errExec, false
			}
			errExec = sqlDriver.Exec(fmt.Sprintf("CREATE DATABASE db%d;", userId))
			if errExec != nil {
				return errExec, false
			}
			errExec = sqlDriver.Exec(fmt.Sprintf("CREATE TABLE [db%d].[dbo].homes (Field INT);", userId))
			if errExec != nil {
				return errExec, false
			}
			actions[userId] = UserExecuteCmd
			cont = true
		}
	}
	if update.Message == nil {
		err = errors.New("message not found")
		return
	}
	if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
		userId = update.Message.From.ID
		if val, ok := actions[userId]; ok {
			switch val {
			default:
				err = errors.New("cmd not found")
				return
			case UserExecuteCmd:
				userScript := update.Message.Text
				errExec := sqlDriver.Exec(userScript)
				if errExec != nil {
					msg := tg.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Во время написания скрипта, вы допустили ошибку...\n [%s]\n Проверьте скрипт и попробуйте еще раз!", errExec.Error()))
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
					return errExec, false
				}
			}
		}
		n := utils.RRange(0, len(TCmds)-1)
		msg := tg.NewMessage(update.Message.Chat.ID, "Вам нужно выполнить случайное задание")
		msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
			tg.NewInlineKeyboardRow(
				tg.NewInlineKeyboardButtonData("Задание", TCmds[n]),
			))
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
	println(userId)
	return
}
