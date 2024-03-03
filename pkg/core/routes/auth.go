package routes

import (
	"errors"
	"fmt"
	"github.com/education-bot/pkg/core"
	"github.com/education-bot/pkg/core/models"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"reflect"
)

const (
	UserRegistrationCmd = "user-reg"
	UserUpdateCmd       = "user-update"
	UserNameCmd         = "user-name-next"
	UserGroupCmd        = "user-group-next"
	UserCourseCmd       = "user-course-next"
)

var (
	usersCmd = map[string]interface{}{
		UserRegistrationCmd: nil,
		UserUpdateCmd:       nil,
		UserNameCmd:         nil,
		UserGroupCmd:        nil,
		UserCourseCmd:       nil,
	}
)

func Authentication(bot *tg.BotAPI, update tg.Update) (err error, cont bool, cUser *core.User) {
	var userId int64
	if update.CallbackQuery != nil {
		userId = update.CallbackQuery.From.ID
		cmd := update.CallbackQuery.Data
		if _, ok := usersCmd[cmd]; !ok {
			cont = true
			return
		}
		switch cmd {
		default:
			err = errors.New("cmd not found")
			return
		case UserRegistrationCmd, UserUpdateCmd:
			msg := tg.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ваша фамилия?")
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			actions[userId] = UserRegistrationCmd
			cont = true
			return
		}
	}
	if update.Message == nil {
		err = errors.New("message not found")
		return
	}
	if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
		userId := update.Message.From.ID
		if val, ok := actions[userId]; ok {
			switch val {
			default:
				return
			case UserRegistrationCmd:
				req := models.UserRequest{
					User: &core.User{
						Id:   userId,
						Name: update.Message.Text,
					},
				}
				_, err = req.Add()
				if err != nil {
					err = errors.New(fmt.Sprintf("что то пошло не так [%s]", err.Error()))
					return
				}
				msg := tg.NewMessage(update.Message.Chat.ID, "Ваша группа?")
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
				actions[userId] = UserNameCmd
			case UserNameCmd:
				user := core.User{Id: userId}
				req := models.UserRequest{}
				getUser, err := req.Get(user.Key())
				if err != nil {
					return errors.New(fmt.Sprintf("что то пошло не так [%s]", err.Error())), false, nil
				}
				getUser.Group = update.Message.Text
				req.User = getUser
				_, err = req.Add()
				if err != nil {
					return errors.New(fmt.Sprintf("что то пошло не так [%s]", err.Error())), false, nil
				}
				msg := tg.NewMessage(update.Message.Chat.ID, "Ваша курс?")
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
					return err, false, nil
				}
				actions[userId] = UserGroupCmd
			case UserGroupCmd:
				user := core.User{Id: userId}
				req := models.UserRequest{}
				getUser, err := req.Get(user.Key())
				if err != nil {
					return errors.New(fmt.Sprintf("что то пошло не так [%s]", err.Error())), false, nil
				}
				getUser.Course = update.Message.Text
				req.User = getUser
				_, err = req.Add()
				if err != nil {
					return errors.New(fmt.Sprintf("что то пошло не так [%s]", err.Error())), false, nil
				}
				msg := tg.NewMessage(update.Message.Chat.ID, "Спасибо за регистрацию")
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
				delete(actions, userId)
			}
		}
		user := &core.User{Id: userId}
		req := models.UserRequest{}
		user, err := req.Get(user.Key())
		if err != nil || user == nil {
			msg := tg.NewMessage(update.Message.Chat.ID, "Вам нужно пройти регистрацию")
			msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
				tg.NewInlineKeyboardRow(
					tg.NewInlineKeyboardButtonData("Регистрация", UserRegistrationCmd),
				))
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			return nil, false, nil
		}
		if user.Group != "" && user.Course != "" && user.Name != "" {
			cUser = user
		}
	}
	return
}
