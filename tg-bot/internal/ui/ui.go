package ui

import (
	"context"
	"log"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"tg-bot/internal/backend"
	"tg-bot/internal/texts"
)

type UI struct {
	be backend.Backend

	mu    sync.Mutex
	langs map[int64]string // language chosen by each user; kept in memory for now
}

func New(be backend.Backend) *UI {
	return &UI{be: be, langs: make(map[int64]string)}
}

func (u *UI) Register(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, u.start)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/cart", bot.MatchTypeExact, u.cartCommand)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "lang:", bot.MatchTypePrefix, u.chooseLang)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "menu:", bot.MatchTypePrefix, u.menu)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "cat:", bot.MatchTypePrefix, u.category)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "prod:", bot.MatchTypePrefix, u.product)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "add:", bot.MatchTypePrefix, u.addToCart)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "size:", bot.MatchTypePrefix, u.sizeGuide)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "cart:", bot.MatchTypePrefix, u.cartAction)
}

func (u *UI) chosenLang(userID int64) (string, bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	lang, ok := u.langs[userID]
	return lang, ok
}

func (u *UI) setLang(userID int64, lang string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.langs[userID] = lang
}

// t returns texts in the user's language: the one they chose,
// otherwise the language of their Telegram app, otherwise Russian.
func (u *UI) t(user models.User) *texts.Texts {
	if lang, ok := u.chosenLang(user.ID); ok {
		return texts.Get(lang)
	}
	return texts.Get(user.LanguageCode)
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
