package bot

import (
	"context"
	"shop_bot/log"
	"time"
)

func (b *Bot) orderDeactivationWorker() {
	t := time.NewTicker(time.Minute)
	for {
		select {
		case <-t.C:
			deactivatedCount, err := b.Storage.DeactivateOrders(context.Background())
			if err != nil {
				log.Error("can't deactivate orders", err)
			} else {
				log.Info("successfully deactivated %d orders", deactivatedCount)
			}
		}
	}
}
