package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/telegram/message/markup"
	"github.com/gotd/td/tg"
	"regexp"
	"shop_bot/log"
	"shop_bot/messages"
	"shop_bot/models"
)

type messageInfo struct {
	e      tg.Entities
	update *tg.UpdateNewMessage
}

func (b *Bot) messageMapping(ctx context.Context, info *messageInfo) error {
	err := b.userInfoCollector(info)
	if err != nil {
		return fmt.Errorf("can't collect user info: %s", err)
	}

	msg := info.update.Message.(*tg.Message)
	text := msg.Message

	if match, _ := regexp.MatchString("/start", text); match {
		return b.StartHandler(ctx, info)
	}
	if match, _ := regexp.MatchString("/help", text); match {
		return b.HelpHandler(ctx, info)
	}
	if match, _ := regexp.MatchString("/catalog", text); match {
		return b.CatalogHandler(ctx, info)
	}
	if match, _ := regexp.MatchString("/myorders", text); match {
		return b.MyOrdersHandler(ctx, info)
	}

	if b.containsImageAndUserIsAdmin(info) {
		return b.downloadImageForLastItem(ctx, info)
	}

	_, err = b.Sender.Answer(info.e, info.update).Text(ctx, messages.UnknownCommand)
	if err != nil {
		return fmt.Errorf("can't reply: %s", err)
	}

	return nil
}

func (b *Bot) userInfoCollector(info *messageInfo) error {
	msg := info.update.Message.(*tg.Message)
	peerUser, ok := msg.PeerID.(*tg.PeerUser)
	if !ok {
		return fmt.Errorf("received message in non-private chat")
	}
	user := info.e.Users[peerUser.UserID]
	log.Info("user: %v", user)

	userModel := &models.User{
		ID:         user.ID,
		ChatID:     user.ID,
		Username:   user.Username,
		AccessHash: &user.AccessHash,
	}

	userFromDB, err := b.Storage.GetUserByID(userModel.ID)
	if err != nil {
		return fmt.Errorf("can't get user to compare: %s", err)
	}
	if userFromDB != nil {
		if userFromDB.AccessHash != userModel.AccessHash ||
			userFromDB.Username != userModel.Username {
			err = b.Storage.UpdateUser(userModel)
			if err != nil {
				return fmt.Errorf("can't update user info: %s", err)
			}
		}
	} else {
		err = b.Storage.CreateUser(userModel)
		if err != nil {
			return fmt.Errorf("can't save user info: %s", err)
		}
	}

	return nil
}

func (b *Bot) StartHandler(ctx context.Context, info *messageInfo) error {
	_, err := b.Sender.Answer(info.e, info.update).Text(ctx, messages.StartReply)
	if err != nil {
		return fmt.Errorf("can't reply: %s", err)
	}
	return nil
}

func (b *Bot) HelpHandler(ctx context.Context, info *messageInfo) error {
	_, err := b.Sender.Answer(info.e, info.update).Text(ctx, messages.HelpReply)
	if err != nil {
		return fmt.Errorf("can't reply: %s", err)
	}
	return nil
}

func (b *Bot) CatalogHandler(ctx context.Context, info *messageInfo) error {
	categories, err := b.Storage.GetTopLevelCategories()
	if err != nil {
		return fmt.Errorf("can't get top level categories: %s", err)
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

	_, err = b.Sender.Answer(info.e, info.update).Markup(markup.InlineKeyboard(rows...)).Text(ctx, messages.CatalogSent)
	if err != nil {
		return fmt.Errorf("can't answer: %s", err)
	}
	return nil
}

func (b *Bot) MyOrdersHandler(ctx context.Context, info *messageInfo) error {
	msg := info.update.Message.(*tg.Message)
	peerUser, ok := msg.PeerID.(*tg.PeerUser)
	if !ok {
		return fmt.Errorf("can't cast peerID to peerUser")
	}
	user := info.e.Users[peerUser.UserID]
	orders, err := b.Storage.GetActiveOrdersByUserID(user.ID)
	if err != nil {
		return fmt.Errorf("can't get active orders for user: %s", err)
	}
	if len(orders) == 0 {
		_, err = b.Sender.Answer(info.e, info.update).Text(ctx, messages.NoActiveOrders)
		if err != nil {
			return fmt.Errorf("can't answer: %s", err)
		}
		return nil
	}
	for _, order := range orders {
		storage, err := b.Storage.GetStorageByID(order.StorageID)
		if err != nil {
			return fmt.Errorf("can't get storage for order: %s", err)
		}
		item, err := b.Storage.GetItemByID(order.ItemID)
		if err != nil {
			return fmt.Errorf("can't get item for order: %s", err)
		}
		msgText := "Вы забронировали %s\n\nВ пункте выдачи %s по адресу %s\n\nКод получения: %d"
		msgText = fmt.Sprintf(msgText, item.Name, storage.Name, storage.Address.String, order.Code)

		rows := []tg.KeyboardButtonRow{{
			Buttons: []tg.KeyboardButtonClass{
				&tg.KeyboardButtonCallback{
					Text: "Отменить бронь",
					Data: []byte(fmt.Sprintf("cancel%d", order.ID)),
				},
			},
		},
		}
		_, err = b.Sender.Answer(info.e, info.update).Markup(markup.InlineKeyboard(rows...)).Text(ctx, msgText)
		if err != nil {
			return fmt.Errorf("can't answer: %s", err)
		}
	}
	return nil
}

func (b *Bot) containsImageAndUserIsAdmin(info *messageInfo) bool {
	msg := info.update.Message.(*tg.Message)
	_, ok := msg.Media.(*tg.MessageMediaPhoto)
	peerUser := msg.PeerID.(*tg.PeerUser)
	var isAdmin bool
	for _, id := range b.Config.Admins {
		if peerUser.UserID == id {
			isAdmin = true
			break
		}
	}
	return ok && isAdmin
}

func (b *Bot) downloadImageForLastItem(ctx context.Context, info *messageInfo) error {
	msg := info.update.Message.(*tg.Message)
	media := msg.Media.(*tg.MessageMediaPhoto)
	photo := media.Photo.(*tg.Photo)
	loc := &tg.InputPhotoFileLocation{
		ID:            photo.ID,
		AccessHash:    photo.AccessHash,
		FileReference: photo.FileReference,
	}
	for i := range photo.Sizes {
		if s, ok := photo.Sizes[i].(*tg.PhotoSizeProgressive); ok {
			loc.ThumbSize = s.Type
		}
	}
	if loc.ThumbSize == "" {
		return errors.New("can't find progressive size for photo")
	}

	loader := downloader.NewDownloader()
	api := b.TelegramClient.API()
	photoData := bytes.NewBuffer(make([]byte, 0))
	fileType, err := loader.Download(api, loc).Stream(ctx, photoData)
	if err != nil {
		return fmt.Errorf("can't download file: %s", err)
	}
	log.Info("image downloaded")

	image := models.Image{
		Bytes: photoData.Bytes(),
	}

	switch fileType.(type) {
	case *tg.StorageFileJpeg:
		image.MimeType = "image/jpeg"
	case *tg.StorageFilePng:
		image.MimeType = "image/png"
	case *tg.StorageFileUnknown:
		return errors.New("unknown image type")
	}

	imageData, err := json.Marshal(image)
	if err != nil {
		return fmt.Errorf("can't marshal image data: %s", err)
	}

	peerUser := msg.PeerID.(*tg.PeerUser)
	userId := peerUser.UserID

	itemId, ok := b.LastItemInChat[userId]
	if !ok {
		_, err = b.Sender.Answer(info.e, info.update).Text(ctx, messages.NewImageSaved)
		if err != nil {
			return fmt.Errorf("can't reply: %s", err)
		}
		return nil
	}
	item, err := b.Storage.GetItemByID(itemId)
	item.Image = imageData
	err = b.Storage.UpdateItem(item)
	if err != nil {
		return fmt.Errorf("can't update item with a new image: %s", err)
	}

	_, err = b.Sender.Answer(info.e, info.update).Text(ctx, messages.NewImageSaved)
	if err != nil {
		return fmt.Errorf("can't reply: %s", err)
	}
	return nil
}
