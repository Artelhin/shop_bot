package bot

import (
	"context"
	"github.com/gotd/td/tg"
)

type callbackInfo struct {
	e      tg.Entities
	update *tg.UpdateBotCallbackQuery
}

func (b *Bot) callbackMapping(ctx context.Context, data string, info *callbackInfo) error {
	return nil
}
