package bot

import (
	"context"
	"fmt"
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
	err := b.userInfoCollector(ctx, info)
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

	_, err = b.Sender.Answer(info.e, info.update).Text(ctx, messages.UnknownCommand)
	if err != nil {
		return fmt.Errorf("can't reply: %s", err)
	}

	return nil
}

func (b *Bot) userInfoCollector(ctx context.Context, info *messageInfo) error {
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
				return fmt.Errorf("can't update user info: %s")
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
