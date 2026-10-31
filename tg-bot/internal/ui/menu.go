package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"tg-bot/internal/texts"
)

func mainMenu(t *texts.Texts) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: t.MenuCatalog, CallbackData: "menu:catalog"},
				{Text: t.MenuCart, CallbackData: "menu:cart"},
			},
			{
				{Text: t.MenuOrders, CallbackData: "menu:orders"},
				{Text: t.MenuContacts, CallbackData: "menu:contacts"},
			},
			{
				{Text: t.MenuInfo, CallbackData: "menu:info"},
				{Text: t.MenuLanguage, CallbackData: "menu:lang"},
			},
		},
	}
}

func langMenu() *models.InlineKeyboardMarkup {
	var row []models.InlineKeyboardButton
	for _, l := range texts.Languages {
		row = append(row, models.InlineKeyboardButton{Text: l.Name, CallbackData: "lang:" + l.Code})
	}
	return &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{row}}
}

func welcome(t *texts.Texts, user models.User) string {
	return fmt.Sprintf(t.Welcome, user.FirstName) + t.DemoBadge
}

// start asks for the language on the first visit, then shows the welcome.
func (u *UI) start(ctx context.Context, b *bot.Bot, update *models.Update) {
	msg := update.Message
	if msg.From == nil {
		return
	}

	params := &bot.SendMessageParams{ChatID: msg.Chat.ID}
	if _, ok := u.chosenLang(msg.From.ID); ok {
		t := u.t(*msg.From)
		params.Text = welcome(t, *msg.From)
		params.ReplyMarkup = mainMenu(t)
	} else {
		params.Text = texts.ChooseLanguage
		params.ReplyMarkup = langMenu()
	}
	_, err := b.SendMessage(ctx, params)
	logErr("send", err)
}

func (u *UI) chooseLang(ctx context.Context, b *bot.Bot, update *models.Update) {
	answer(ctx, b, update, "")
	msg := callbackMessage(update)
	if msg == nil {
		return
	}

	user := update.CallbackQuery.From
	lang := strings.TrimPrefix(update.CallbackQuery.Data, "lang:")
	if texts.Supported(lang) {
		u.setLang(user.ID, lang)
	}
	t := u.t(user)
	showText(ctx, b, msg, welcome(t, user), mainMenu(t))
}

func (u *UI) menu(ctx context.Context, b *bot.Bot, update *models.Update) {
	answer(ctx, b, update, "")
	msg := callbackMessage(update)
	if msg == nil {
		return
	}
	user := update.CallbackQuery.From
	t := u.t(user)

	var text string
	switch update.CallbackQuery.Data {
	case "menu:home":
		showText(ctx, b, msg, welcome(t, user), mainMenu(t))
		return
	case "menu:catalog":
		u.showCategories(ctx, b, msg, t)
		return
	case "menu:lang":
		showText(ctx, b, msg, texts.ChooseLanguage, langMenu())
		return
	case "menu:cart":
		text, kb := u.cartScreen(ctx, user.ID, t)
		showText(ctx, b, msg, text, kb)
		return
	case "menu:orders":
		text = t.NoOrders
	case "menu:contacts":
		text = t.ContactsSoon
	case "menu:info":
		text = t.Info
	default:
		text = t.UnknownButton
	}
	showText(ctx, b, msg, text+t.DemoBadge, mainMenu(t))
}
