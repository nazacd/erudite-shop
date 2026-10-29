package texts

import (
	"strconv"
	"strings"
)

const (
	DemoBadge = "\n🧪 Демо-версия — оплата не списывается"

	Welcome       = "👋 Добро пожаловать в «Демо-Магазин»!\nВыберите раздел"
	About         = "ℹ️ «Демо-Магазин» — пример бота для онлайн-торговли."
	CartEmpty     = "🛒 Корзина пуста."
	NoOrders      = "📦 У вас пока нет заказов."
	ContactsSoon  = "📞 Контакты — скоро."
	UnknownButton = "Эта кнопка устарела. Откройте меню заново."

	CatalogTitle   = "🛍 Каталог\nВыберите категорию:"
	ChooseProduct  = "Выберите товар:"
	EmptyCategory  = "В этой категории пока нет товаров."
	ProductGone    = "Этот товар больше недоступен."
	SomethingWrong = "Что-то пошло не так. Попробуйте ещё раз."

	Back          = "⬅️ Назад"
	AddToCart     = "🛒 В корзину"
	AddToCartSoon = "Корзина появится в следующем уроке 🙂"
)

func FormatSum(n int64) string {
	s := strconv.FormatInt(n, 10)
	var sb strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			sb.WriteString("\u00a0")
		}
		sb.WriteRune(r)
	}
	return sb.String() + "\u00a0сум"
}
