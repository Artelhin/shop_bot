package bot

import (
	"context"
	"fmt"
	"github.com/gotd/td/tg"
	"shop_bot/log"
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

func (b *Bot) StartHandler() {}
