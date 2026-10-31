package texts

import (
	"strconv"
	"strings"
)

// ChooseLanguage is shown before we know the user's language, so it is in all three.
const ChooseLanguage = "🌐 Tilni tanlang / Выберите язык / Choose a language"

// Languages are the language buttons, in the order they are shown.
var Languages = []struct{ Code, Name string }{
	{"uz", "🇺🇿 O‘zbekcha"},
	{"ru", "🇷🇺 Русский"},
	{"en", "🇬🇧 English"},
}

type Texts struct {
	DemoBadge string

	Welcome       string // %s is the user's first name
	Info          string
	CartEmpty     string
	NoOrders      string
	ContactsSoon  string
	UnknownButton string

	MenuCatalog  string
	MenuCart     string
	MenuOrders   string
	MenuContacts string
	MenuInfo     string
	MenuLanguage string

	CatalogTitle   string
	ChooseProduct  string
	EmptyCategory  string
	ProductGone    string
	SomethingWrong string

	SizeGuideButton string
	SizeGuide       string

	Back      string
	AddToCart string
	Added     string

	CartTitle    string
	CartTotal    string // %s is the formatted sum
	ToCatalog    string
	ClearCart    string
	Checkout     string
	CheckoutSoon string

	Currency string
}

var uz = Texts{
	DemoBadge: "\n🧪 Demo versiya — to‘lov olinmaydi",

	Welcome: "👋 Salom, %s! Erudite Shop’ga xush kelibsiz — Erudite Education maktabining rasmiy kiyim do‘koni.\n\n" +
		"Bu yerda siz:\n" +
		"👕 Xudi, futbolka va boshqa kiyimlarni ko‘rishingiz\n" +
		"🛒 Bir necha bosishda buyurtma berishingiz\n" +
		"📦 Buyurtmangizni kuzatishingiz mumkin\n\n" +
		"Boshlash uchun quyidagi bo‘limni tanlang 👇",
	Info: "ℹ️ Yetkazib berish va to‘lov\n\n" +
		"🏫 Olib ketish — Erudite Education maktabidan\n" +
		"🚚 Yetkazib berish — Toshkent bo‘ylab\n" +
		"💳 To‘lov — Click, Payme yoki naqd pul",
	CartEmpty:     "🛒 Savat bo‘sh.",
	NoOrders:      "📦 Sizda hali buyurtmalar yo‘q.",
	ContactsSoon:  "📞 Kontaktlar — tez orada.",
	UnknownButton: "Bu tugma eskirgan. Menyuni qaytadan oching.",

	MenuCatalog:  "🛍 Katalog",
	MenuCart:     "🛒 Savat",
	MenuOrders:   "📦 Buyurtmalarim",
	MenuContacts: "📞 Biz bilan bog‘lanish",
	MenuInfo:     "ℹ️ Ma’lumot",
	MenuLanguage: "🌐 Til",

	CatalogTitle:   "🛍 Katalog\nKategoriyani tanlang:",
	ChooseProduct:  "Mahsulotni tanlang:",
	EmptyCategory:  "Bu kategoriyada hali mahsulotlar yo‘q.",
	ProductGone:    "Bu mahsulot endi mavjud emas.",
	SomethingWrong: "Nimadir xato ketdi. Qaytadan urinib ko‘ring.",

	SizeGuideButton: "📏 O‘lchamlar jadvali",
	SizeGuide: "📏 O‘lchamlar jadvali (ko‘krak aylanasi)\n\n" +
		"XS — 82–86 sm\nS — 86–90 sm\nM — 90–96 sm\nL — 96–102 sm\nXL — 102–108 sm\nXXL — 108–114 sm\n\n" +
		"Ko‘krakni eng keng joyidan o‘lchang. Ikki o‘lcham orasida bo‘lsangiz, kattarog‘ini tanlang.",

	Back:      "⬅️ Orqaga",
	AddToCart: "🛒 Savatga",
	Added:     "Qo‘shildi ✅",

	CartTitle:    "🛒 Savat",
	CartTotal:    "Jami: %s",
	ToCatalog:    "🛍 Katalogga",
	ClearCart:    "🗑 Tozalash",
	Checkout:     "✅ Buyurtma berish",
	CheckoutSoon: "Buyurtma berish keyingi darsda bo‘ladi",

	Currency: "so‘m",
}

var ru = Texts{
	DemoBadge: "\n🧪 Демо-версия — оплата не списывается",

	Welcome: "👋 Привет, %s! Добро пожаловать в Erudite Shop — официальный магазин одежды школы Erudite Education.\n\n" +
		"Здесь вы можете:\n" +
		"👕 Смотреть худи, футболки и многое другое\n" +
		"🛒 Оформить заказ в пару нажатий\n" +
		"📦 Отслеживать свой заказ\n\n" +
		"Выберите раздел ниже 👇",
	Info: "ℹ️ Доставка и оплата\n\n" +
		"🏫 Самовывоз — в школе Erudite Education\n" +
		"🚚 Доставка — по Ташкенту\n" +
		"💳 Оплата — Click, Payme или наличными",
	CartEmpty:     "🛒 Корзина пуста.",
	NoOrders:      "📦 У вас пока нет заказов.",
	ContactsSoon:  "📞 Контакты — скоро.",
	UnknownButton: "Эта кнопка устарела. Откройте меню заново.",

	MenuCatalog:  "🛍 Каталог",
	MenuCart:     "🛒 Корзина",
	MenuOrders:   "📦 Мои заказы",
	MenuContacts: "📞 Связаться с нами",
	MenuInfo:     "ℹ️ Инфо",
	MenuLanguage: "🌐 Язык",

	CatalogTitle:   "🛍 Каталог\nВыберите категорию:",
	ChooseProduct:  "Выберите товар:",
	EmptyCategory:  "В этой категории пока нет товаров.",
	ProductGone:    "Этот товар больше недоступен.",
	SomethingWrong: "Что-то пошло не так. Попробуйте ещё раз.",

	SizeGuideButton: "📏 Таблица размеров",
	SizeGuide: "📏 Таблица размеров (обхват груди)\n\n" +
		"XS — 82–86 см\nS — 86–90 см\nM — 90–96 см\nL — 96–102 см\nXL — 102–108 см\nXXL — 108–114 см\n\n" +
		"Измеряйте грудь по самой широкой части. Если вы между размерами, берите больший.",

	Back:      "⬅️ Назад",
	AddToCart: "🛒 В корзину",
	Added:     "Добавлено ✅",

	CartTitle:    "🛒 Корзина",
	CartTotal:    "Итого: %s",
	ToCatalog:    "🛍 В каталог",
	ClearCart:    "🗑 Очистить",
	Checkout:     "✅ Оформить заказ",
	CheckoutSoon: "Оформление будет в следующем уроке",

	Currency: "сум",
}

var en = Texts{
	DemoBadge: "\n🧪 Demo version — no payment is taken",

	Welcome: "👋 Hi, %s! Welcome to Erudite Shop, the official clothing store of Erudite Education School.\n\n" +
		"Here you can:\n" +
		"👕 Browse hoodies, T-shirts, and more\n" +
		"🛒 Order in a few taps\n" +
		"📦 Track your order\n\n" +
		"Choose an option below to get started 👇",
	Info: "ℹ️ Delivery and payment\n\n" +
		"🏫 Pickup — at Erudite Education School\n" +
		"🚚 Delivery — within Tashkent\n" +
		"💳 Payment — Click, Payme, or cash",
	CartEmpty:     "🛒 Your cart is empty.",
	NoOrders:      "📦 You have no orders yet.",
	ContactsSoon:  "📞 Contacts — coming soon.",
	UnknownButton: "This button is outdated. Please open the menu again.",

	MenuCatalog:  "🛍 Catalog",
	MenuCart:     "🛒 Cart",
	MenuOrders:   "📦 My orders",
	MenuContacts: "📞 Contact us",
	MenuInfo:     "ℹ️ Info",
	MenuLanguage: "🌐 Language",

	CatalogTitle:   "🛍 Catalog\nChoose a category:",
	ChooseProduct:  "Choose a product:",
	EmptyCategory:  "There are no products in this category yet.",
	ProductGone:    "This product is no longer available.",
	SomethingWrong: "Something went wrong. Please try again.",

	SizeGuideButton: "📏 Size guide",
	SizeGuide: "📏 Size guide (chest)\n\n" +
		"XS — 82–86 cm\nS — 86–90 cm\nM — 90–96 cm\nL — 96–102 cm\nXL — 102–108 cm\nXXL — 108–114 cm\n\n" +
		"Measure around the widest part of the chest. Between two sizes? Take the larger one.",

	Back:      "⬅️ Back",
	AddToCart: "🛒 Add to cart",
	Added:     "Added ✅",

	CartTitle:    "🛒 Cart",
	CartTotal:    "Total: %s",
	ToCatalog:    "🛍 To catalog",
	ClearCart:    "🗑 Clear",
	Checkout:     "✅ Checkout",
	CheckoutSoon: "Checkout is coming in the next lesson",

	Currency: "UZS",
}

// Supported reports whether lang is one of the bot's languages.
func Supported(lang string) bool {
	return lang == "uz" || lang == "ru" || lang == "en"
}

// Get returns the texts for lang, falling back to Russian.
func Get(lang string) *Texts {
	switch lang {
	case "uz":
		return &uz
	case "en":
		return &en
	default:
		return &ru
	}
}

func (t *Texts) FormatSum(n int64) string {
	s := strconv.FormatInt(n, 10)
	var sb strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			sb.WriteString(" ")
		}
		sb.WriteRune(r)
	}
	return sb.String() + " " + t.Currency
}
