package ui

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"shopbot/internal/backend"
)

type UI struct {
	be backend.Backend
}

func New(be backend.Backend) *UI {
	return &UI{be: be}
}

func (u *UI) Register(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, u.start)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "menu:", bot.MatchTypePrefix, u.menu)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "cat:", bot.MatchTypePrefix, u.category)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "prod:", bot.MatchTypePrefix, u.product)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "add:", bot.MatchTypePrefix, u.addToCart)
}

func logErr(what string, err error) {
	if err != nil {
		log.Println(what, "error:", err)
	}
}

func answer(ctx context.Context, b *bot.Bot, update *models.Update, text string) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            text,
	})
	logErr("answer", err)
}

func callbackMessage(update *models.Update) *models.Message {
	if update.CallbackQuery == nil {
		return nil
	}
	return update.CallbackQuery.Message.Message
}

func showText(ctx context.Context, b *bot.Bot, msg *models.Message, text string, kb *models.InlineKeyboardMarkup) {
	if len(msg.Photo) > 0 {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      msg.Chat.ID,
			Text:        text,
			ReplyMarkup: kb,
		})
		logErr("send", err)
		_, err = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    msg.Chat.ID,
			MessageID: msg.ID,
		})
		logErr("delete", err)
		return
	}
	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      msg.Chat.ID,
		MessageID:   msg.ID,
		Text:        text,
		ReplyMarkup: kb,
	})
	logErr("edit", err)
}
