package ui

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"tg-bot/internal/backend"
	"tg-bot/internal/texts"
)

// backTo is a keyboard with one "Назад" button.
func backTo(data string) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: texts.Back, CallbackData: data}},
		},
	}
}

// showCategories: the list of categories, one button per row.
func (u *UI) showCategories(ctx context.Context, b *bot.Bot, msg *models.Message) {
	cats, err := u.be.Categories(ctx)
	if err != nil {
		logErr("categories", err)
		showText(ctx, b, msg, texts.SomethingWrong+texts.DemoBadge, backTo("menu:home"))
		return
	}

	var rows [][]models.InlineKeyboardButton
	for _, c := range cats {
		rows = append(rows, []models.InlineKeyboardButton{
			{Text: c.Name, CallbackData: "cat:" + c.ID},
		})
	}
	rows = append(rows, []models.InlineKeyboardButton{
		{Text: texts.Back, CallbackData: "menu:home"},
	})

	showText(ctx, b, msg, texts.CatalogTitle+texts.DemoBadge,
		&models.InlineKeyboardMarkup{InlineKeyboard: rows})
}

// category: the products of one category.
func (u *UI) category(ctx context.Context, b *bot.Bot, update *models.Update) {
	answer(ctx, b, update, "")
	msg := callbackMessage(update)
	if msg == nil {
		return
	}
	categoryID := strings.TrimPrefix(update.CallbackQuery.Data, "cat:")

	products, err := u.be.ProductsByCategory(ctx, categoryID)
	if err != nil {
		logErr("products", err)
		showText(ctx, b, msg, texts.SomethingWrong+texts.DemoBadge, backTo("menu:catalog"))
		return
	}
	if len(products) == 0 {
		showText(ctx, b, msg, texts.EmptyCategory+texts.DemoBadge, backTo("menu:catalog"))
		return
	}

	var rows [][]models.InlineKeyboardButton
	for _, p := range products {
		rows = append(rows, []models.InlineKeyboardButton{
			{Text: p.Name + " — " + texts.FormatSum(p.Price), CallbackData: "prod:" + p.ID},
		})
	}
	rows = append(rows, []models.InlineKeyboardButton{
		{Text: texts.Back, CallbackData: "menu:catalog"},
	})

	showText(ctx, b, msg, texts.ChooseProduct+texts.DemoBadge,
		&models.InlineKeyboardMarkup{InlineKeyboard: rows})
}

// product: the card with a photo.
func (u *UI) product(ctx context.Context, b *bot.Bot, update *models.Update) {
	answer(ctx, b, update, "")
	msg := callbackMessage(update)
	if msg == nil {
		return
	}
	productID := strings.TrimPrefix(update.CallbackQuery.Data, "prod:")

	p, err := u.be.Product(ctx, productID)
	if errors.Is(err, backend.ErrNotFound) {
		showText(ctx, b, msg, texts.ProductGone+texts.DemoBadge, backTo("menu:catalog"))
		return
	}
	if err != nil {
		logErr("product", err)
		showText(ctx, b, msg, texts.SomethingWrong+texts.DemoBadge, backTo("menu:catalog"))
		return
	}

	caption := fmt.Sprintf("%s\n\n%s\n\n💰 %s", p.Name, p.Description, texts.FormatSum(p.Price)) + texts.DemoBadge
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: texts.AddToCart, CallbackData: "add:" + p.ID}},
			{{Text: texts.Back, CallbackData: "cat:" + p.CategoryID}},
		},
	}

	_, err = b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:      msg.Chat.ID,
		Photo:       &models.InputFileString{Data: p.PhotoURL},
		Caption:     caption,
		ReplyMarkup: kb,
	})
	if err != nil {
		// The photo could not be loaded: show the same card as text instead.
		log.Println("photo error:", err)
		showText(ctx, b, msg, caption, kb)
		return
	}
	_, err = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    msg.Chat.ID,
		MessageID: msg.ID,
	})
	logErr("delete", err)
}

// addToCart is a placeholder until Lesson 4.
func (u *UI) addToCart(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            texts.AddToCartSoon,
		ShowAlert:       true,
	})
	logErr("answer", err)
}
