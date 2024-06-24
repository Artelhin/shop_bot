package messages

const (
	StartReply = `Привет! Здесь ты можешь забронировать девайсы из каталога. В любой момент ты можешь вызвать /catalog, чтобы посмотреть полный список товаров`
	HelpReply  = `/catalog - список товаров
/myorders - посмотреть список активных заказов`
	UnknownCommand  = `Неизвестная команда`
	CatalogSent     = `Предоставляю каталог`
	NotInStock      = `К сожалению товара уже нет в наличии, попробуйте другой склад или повторите заказ позже`
	NoActiveOrders  = `У вас нет активных бронирований! Возможно, вам понравится что-то из каталога /catalog ?`
	OrderCancelled  = `Бронь успешно отменена!`
	NewImageSaved   = `Новое фото товара сохранено!`
	ItemNotSelected = `Товар для обновления не выбран, пожалуйста, выберете товар из каталога!`
)
