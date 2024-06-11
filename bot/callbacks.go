package bot

import (
	"context"
	"fmt"
	"github.com/gotd/td/telegram/message/markup"
	"github.com/gotd/td/tg"
	"math/rand"
	"regexp"
	"shop_bot/messages"
	"shop_bot/models"
	"shop_bot/storage"
	"strconv"
	"strings"

	"shop_bot/log"
)

type callbackInfo struct {
	e      tg.Entities
	update *tg.UpdateBotCallbackQuery
}

func (b *Bot) callbackMapping(ctx context.Context, data string, info *callbackInfo) error {
	log.Info("got callback data=%s", data)

	_, err := b.TelegramAPI.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
		QueryID: info.update.QueryID,
	})
	if err != nil {
		log.Error("can't send answer to callback request", err)
	}

	if match, _ := regexp.MatchString("showcat*", data); match {
		return b.showCategoryCallback(ctx, data, info)
	}
	if match, _ := regexp.MatchString("showitem*", data); match {
		return b.showItemCallback(ctx, data, info)
	}
	if match, _ := regexp.MatchString("order*", data); match {
		return b.orderCallback(ctx, data, info)
	}

	return nil
}

func (b *Bot) showCategoryCallback(ctx context.Context, data string, info *callbackInfo) error {
	strid := strings.TrimPrefix(data, "showcat")
	id, err := strconv.ParseUint(strid, 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse id for category: %s", err)
	}
	category, err := b.Storage.GetCategoryByID(int64(id))
	if err != nil {
		return fmt.Errorf("can't get category: %s", err)
	}

	categories, err := b.Storage.GetSubcategoriesByCategoryID(int64(id))
	if err != nil {
		return fmt.Errorf("can't get sub categories: %s", err)
	}

	items, err := b.Storage.GetItemsByCategoryID(int64(id))
	if err != nil {
		return fmt.Errorf("can't get items for category: %s", err)
	}

	rows := make([]tg.KeyboardButtonRow, 0)
	for i := range categories {
		rows = append(rows, tg.KeyboardButtonRow{
			Buttons: []tg.KeyboardButtonClass{
				&tg.KeyboardButtonCallback{
					Text: categories[i].Name,
					Data: []byte(fmt.Sprintf("showcat%d", categories[i].ID)),
				},
			},
		})
	}
	for i := range items {
		rows = append(rows, tg.KeyboardButtonRow{
			Buttons: []tg.KeyboardButtonClass{
				&tg.KeyboardButtonCallback{
					Text: items[i].Name,
					Data: []byte(fmt.Sprintf("showitem%d", items[i].ID)),
				},
			},
		})
	}
	if category.ParentID != nil {
		rows = append(rows, tg.KeyboardButtonRow{
			Buttons: []tg.KeyboardButtonClass{
				&tg.KeyboardButtonCallback{
					Text: "Назад",
					Data: []byte(fmt.Sprintf("showcat%d", *category.ParentID)),
				},
			},
		})
	}

	user, err := b.Storage.GetUserByID(info.update.UserID)
	if err != nil {
		return fmt.Errorf("can't get user from storage: %s", err)
	}
	peerUser := &tg.InputPeerUser{
		UserID:     user.ID,
		AccessHash: *user.AccessHash,
	}
	_, err = b.Sender.To(peerUser).Markup(markup.InlineKeyboard(rows...)).Edit(info.update.MsgID).Text(ctx, category.Name)
	if err != nil {
		return fmt.Errorf("can't answer: %s", err)
	}
	return nil
}

func (b *Bot) showItemCallback(ctx context.Context, data string, info *callbackInfo) error {
	strid := strings.TrimPrefix(data, "showitem")
	id, err := strconv.ParseUint(strid, 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse id for item: %s", err)
	}
	item, err := b.Storage.GetItemByID(int64(id))
	if err != nil {
		return fmt.Errorf("can't get item: %s", err)
	}

	storages, err := b.Storage.GetStoragesForItemID(item.ID)
	if err != nil {
		return fmt.Errorf("can't get storages for item: %s", err)
	}

	msgParts := []string{item.Name}
	if item.Description.Valid && item.Description.String != "" {
		msgParts = append(msgParts, item.Description.String)
	}
	if len(storages) == 0 {
		msgParts = append(msgParts, "\nНет в наличии")
	} else {
		msgParts = append(msgParts, "\nЕсть в наличии на складах:")
		for i, store := range storages {
			if store.Address.Valid && store.Address.String != "" {
				msgParts = append(msgParts, fmt.Sprintf("%d) %s (%s)", i+1, store.Name, store.Address.String))
			} else {
				msgParts = append(msgParts, fmt.Sprintf("%d) %s", i, store.Name))
			}
		}
	}

	rows := make([]tg.KeyboardButtonRow, 0)
	for i := range storages {
		rows = append(rows, tg.KeyboardButtonRow{
			Buttons: []tg.KeyboardButtonClass{
				&tg.KeyboardButtonCallback{
					Text: fmt.Sprintf("Забронировать в %s", storages[i].Name),
					Data: []byte(fmt.Sprintf("order%dfrom%d", item.ID, storages[i].ID)),
				},
			},
		})
	}

	user, err := b.Storage.GetUserByID(info.update.UserID)
	if err != nil {
		return fmt.Errorf("can't get user from storage: %s", err)
	}
	peerUser := &tg.InputPeerUser{
		UserID:     user.ID,
		AccessHash: *user.AccessHash,
	}

	if len(rows) == 0 {
		_, err = b.Sender.To(peerUser).Text(ctx, strings.Join(msgParts, "\n\n"))
		if err != nil {
			return fmt.Errorf("can't answer: %s", err)
		}
	} else {
		_, err = b.Sender.To(peerUser).Markup(markup.InlineKeyboard(rows...)).Text(ctx, strings.Join(msgParts, "\n"))
		if err != nil {
			return fmt.Errorf("can't answer: %s", err)
		}
	}

	return nil
}

func (b *Bot) orderCallback(ctx context.Context, data string, info *callbackInfo) error {
	data = strings.TrimPrefix(data, "order")
	ids := strings.Split(data, "from")

	itemId, err := strconv.ParseUint(ids[0], 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse id for item: %s", err)
	}
	storageId, err := strconv.ParseUint(ids[0], 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse id for storage: %s", err)
	}
	store, err := b.Storage.GetStorageByID(int64(storageId))
	if err != nil {
		return fmt.Errorf("can't get storage: %s", err)
	}

	user, err := b.Storage.GetUserByID(info.update.UserID)
	if err != nil {
		return fmt.Errorf("can't get user from storage: %s", err)
	}

	code := rand.Int63()

	order := &models.Order{
		UserID:    user.ID,
		ItemID:    int64(itemId),
		StorageID: int64(storageId),
		Active:    true,
		Code:      code,
	}
	result, err := b.Storage.CreateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("can't create new order: %s", err)
	}

	peerUser := &tg.InputPeerUser{
		UserID:     user.ID,
		AccessHash: *user.AccessHash,
	}

	if result == storage.OrderResultNotInStock {
		_, err = b.Sender.To(peerUser).Text(ctx, messages.NotInStock)
		if err != nil {
			return fmt.Errorf("can't send message: %s", err)
		}
		return nil
	}
	if result == storage.OrderResultSuccess {
		msg := fmt.Sprintf("Ваш заказ будет ждать вас на складе %s по адресу %s в течение суток. Код для получения - %d. Спасибо за использование нашего бота!",
			store.Name, store.Address.String, order.Code)
		_, err = b.Sender.To(peerUser).Text(ctx, msg)
		if err != nil {
			return fmt.Errorf("can't send message: %s", err)
		}
	}
	return nil
}
