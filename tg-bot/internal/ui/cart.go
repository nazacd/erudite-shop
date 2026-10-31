package ui

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"tg-bot/internal/backend"
	"tg-bot/internal/texts"
)

// cartScreen builds the text and keyboard of the user's cart.
// It only builds them; the caller decides whether to send a new message or edit one.
func (u *UI) cartScreen(ctx context.Context, userID int64, t *texts.Texts) (string, *models.InlineKeyboardMarkup) {
	cart, err := u.be.Cart(ctx, userID)
	if err != nil {
		logErr("cart", err)
		return t.SomethingWrong + t.DemoBadge, backTo(t, "menu:home")
	}

	if len(cart.Items) == 0 {
		return t.CartEmpty + t.DemoBadge, &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{{Text: t.ToCatalog, CallbackData: "menu:catalog"}},
				{{Text: t.Back, CallbackData: "menu:home"}},
			},
		}
	}

	var sb strings.Builder
	sb.WriteString(t.CartTitle + "\n\n")
	var rows [][]models.InlineKeyboardButton
	for _, item := range cart.Items {
		p := item.Product
		// "Футболка Classic — 2 × 120 000 сум"
		fmt.Fprintf(&sb, "%s — %d × %s\n", p.Name, item.Qty, t.FormatSum(p.Price))
		rows = append(rows, []models.InlineKeyboardButton{
			{Text: "➖", CallbackData: "cart:dec:" + p.ID},
			{Text: strconv.Itoa(item.Qty), CallbackData: "cart:noop"},
			{Text: "➕", CallbackData: "cart:inc:" + p.ID},
			{Text: "❌", CallbackData: "cart:del:" + p.ID},
		})
	}
	sb.WriteString("\n" + fmt.Sprintf(t.CartTotal, t.FormatSum(cart.Total)))
	sb.WriteString(t.DemoBadge)

	rows = append(rows,
		[]models.InlineKeyboardButton{{Text: t.ClearCart, CallbackData: "cart:clear"}},
		[]models.InlineKeyboardButton{{Text: t.Checkout, CallbackData: "cart:checkout"}},
		[]models.InlineKeyboardButton{{Text: t.Back, CallbackData: "menu:home"}},
	)
	return sb.String(), &models.InlineKeyboardMarkup{InlineKeyboard: rows}
}

// cartCommand: /cart sends the cart as a new message.
func (u *UI) cartCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	msg := update.Message
	if msg.From == nil {
		return
	}
	text, kb := u.cartScreen(ctx, msg.From.ID, u.t(*msg.From))
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      msg.Chat.ID,
		Text:        text,
		ReplyMarkup: kb,
	})
	logErr("send", err)
}

// cartAction handles the cart buttons: "cart:inc:p1", "cart:dec:p1", "cart:del:p1",
// "cart:clear", "cart:checkout" and "cart:noop". After a change it redraws the cart
// in the same message.
func (u *UI) cartAction(ctx context.Context, b *bot.Bot, update *models.Update) {
	msg := callbackMessage(update)
	if msg == nil {
		answer(ctx, b, update, "")
		return
	}
	user := update.CallbackQuery.From
	t := u.t(user)

	// "cart:inc:p1" → ["cart", "inc", "p1"]; "cart:clear" → ["cart", "clear"]
	parts := strings.Split(update.CallbackQuery.Data, ":")
	action := parts[1]
	productID := ""
	if len(parts) > 2 {
		productID = parts[2]
	}

	var err error
	switch action {
	case "noop":
		answer(ctx, b, update, "")
		return
	case "checkout":
		answer(ctx, b, update, t.CheckoutSoon)
		return
	case "inc":
		err = u.be.ChangeQty(ctx, user.ID, productID, +1)
	case "dec":
		err = u.be.ChangeQty(ctx, user.ID, productID, -1)
	case "del":
		err = u.be.RemoveFromCart(ctx, user.ID, productID)
	case "clear":
		err = u.be.ClearCart(ctx, user.ID)
	default:
		answer(ctx, b, update, t.UnknownButton)
		return
	}

	// ErrNotFound means the item was already gone (a button on an old message).
	// That is not a problem: we just redraw the cart as it is now.
	if err != nil && !errors.Is(err, backend.ErrNotFound) {
		logErr("cart "+action, err)
		answer(ctx, b, update, t.SomethingWrong)
		return
	}
	answer(ctx, b, update, "")

	text, kb := u.cartScreen(ctx, user.ID, t)
	showText(ctx, b, msg, text, kb)
}
