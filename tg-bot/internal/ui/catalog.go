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

// backTo is a keyboard with one "Back" button.
func backTo(t *texts.Texts, data string) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: t.Back, CallbackData: data}},
		},
	}
}

// showCategories: the list of categories, one button per row.
func (u *UI) showCategories(ctx context.Context, b *bot.Bot, msg *models.Message, t *texts.Texts) {
	cats, err := u.be.Categories(ctx)
	if err != nil {
		logErr("categories", err)
		showText(ctx, b, msg, t.SomethingWrong+t.DemoBadge, backTo(t, "menu:home"))
		return
	}

	var rows [][]models.InlineKeyboardButton
	for _, c := range cats {
		rows = append(rows, []models.InlineKeyboardButton{
			{Text: c.Name, CallbackData: "cat:" + c.ID},
		})
	}
	rows = append(rows,
		[]models.InlineKeyboardButton{{Text: t.SizeGuideButton, CallbackData: "size:menu:catalog"}},
		[]models.InlineKeyboardButton{{Text: t.Back, CallbackData: "menu:home"}},
	)

	showText(ctx, b, msg, t.CatalogTitle+t.DemoBadge,
		&models.InlineKeyboardMarkup{InlineKeyboard: rows})
}

// category: the products of one category.
func (u *UI) category(ctx context.Context, b *bot.Bot, update *models.Update) {
	answer(ctx, b, update, "")
	msg := callbackMessage(update)
	if msg == nil {
		return
	}
	t := u.t(update.CallbackQuery.From)
	categoryID := strings.TrimPrefix(update.CallbackQuery.Data, "cat:")

	products, err := u.be.ProductsByCategory(ctx, categoryID)
	if err != nil {
		logErr("products", err)
		showText(ctx, b, msg, t.SomethingWrong+t.DemoBadge, backTo(t, "menu:catalog"))
		return
	}
	if len(products) == 0 {
		showText(ctx, b, msg, t.EmptyCategory+t.DemoBadge, backTo(t, "menu:catalog"))
		return
	}

	var rows [][]models.InlineKeyboardButton
	for _, p := range products {
		rows = append(rows, []models.InlineKeyboardButton{
			{Text: p.Name + " — " + t.FormatSum(p.Price), CallbackData: "prod:" + p.ID},
		})
	}
	rows = append(rows, []models.InlineKeyboardButton{
		{Text: t.Back, CallbackData: "menu:catalog"},
	})

	showText(ctx, b, msg, t.ChooseProduct+t.DemoBadge,
		&models.InlineKeyboardMarkup{InlineKeyboard: rows})
}

// product: the card with a photo.
func (u *UI) product(ctx context.Context, b *bot.Bot, update *models.Update) {
	answer(ctx, b, update, "")
	msg := callbackMessage(update)
	if msg == nil {
		return
	}
	t := u.t(update.CallbackQuery.From)
	productID := strings.TrimPrefix(update.CallbackQuery.Data, "prod:")

	p, err := u.be.Product(ctx, productID)
	if errors.Is(err, backend.ErrNotFound) {
		showText(ctx, b, msg, t.ProductGone+t.DemoBadge, backTo(t, "menu:catalog"))
		return
	}
	if err != nil {
		logErr("product", err)
		showText(ctx, b, msg, t.SomethingWrong+t.DemoBadge, backTo(t, "menu:catalog"))
		return
	}

	caption := fmt.Sprintf("%s\n\n%s\n\n💰 %s", p.Name, p.Description, t.FormatSum(p.Price)) + t.DemoBadge
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: t.AddToCart, CallbackData: "add:" + p.ID}},
			{{Text: t.Back, CallbackData: "cat:" + p.CategoryID}},
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

// addToCart puts the product into the user's cart and confirms with a short pop-up.
func (u *UI) addToCart(ctx context.Context, b *bot.Bot, update *models.Update) {
	user := update.CallbackQuery.From
	t := u.t(user)
	productID := strings.TrimPrefix(update.CallbackQuery.Data, "add:")

	err := u.be.AddToCart(ctx, user.ID, productID)
	switch {
	case errors.Is(err, backend.ErrNotFound):
		answer(ctx, b, update, t.ProductGone)
	case err != nil:
		logErr("add to cart", err)
		answer(ctx, b, update, t.SomethingWrong)
	default:
		answer(ctx, b, update, t.Added)
	}
}

// sizeGuide shows the size table. The callback data after "size:"
// is where the Back button returns, e.g. "size:menu:catalog".
func (u *UI) sizeGuide(ctx context.Context, b *bot.Bot, update *models.Update) {
	answer(ctx, b, update, "")
	msg := callbackMessage(update)
	if msg == nil {
		return
	}
	t := u.t(update.CallbackQuery.From)
	back := strings.TrimPrefix(update.CallbackQuery.Data, "size:")
	showText(ctx, b, msg, t.SizeGuide+t.DemoBadge, backTo(t, back))
}
