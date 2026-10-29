package ui

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"tg-bot/internal/texts"
)

func mainMenu() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "🛍 Каталог", CallbackData: "menu:catalog"},
				{Text: "🛒 Корзина", CallbackData: "menu:cart"},
			},
			{
				{Text: "📦 Мои заказы", CallbackData: "menu:orders"},
				{Text: "📞 Контакты", CallbackData: "menu:contacts"},
			},
			{
				{Text: "ℹ️ О магазине", CallbackData: "menu:about"},
			},
		},
	}
}

func (u *UI) start(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        texts.Welcome + texts.DemoBadge,
		ReplyMarkup: mainMenu(),
	})
	logErr("send", err)
}

func (u *UI) menu(ctx context.Context, b *bot.Bot, update *models.Update) {
	answer(ctx, b, update, "")
	msg := callbackMessage(update)
	if msg == nil {
		return
	}

	var text string
	switch update.CallbackQuery.Data {
	case "menu:home":
		text = texts.Welcome
	case "menu:catalog":
		u.showCategories(ctx, b, msg)
		return
	case "menu:cart":
		text = texts.CartEmpty
	case "menu:orders":
		text = texts.NoOrders
	case "menu:contacts":
		text = texts.ContactsSoon
	case "menu:about":
		text = texts.About
	default:
		text = texts.UnknownButton
	}
	showText(ctx, b, msg, text+texts.DemoBadge, mainMenu())
}
