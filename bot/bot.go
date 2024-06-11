package bot

import (
	"context"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"os"
	"shop_bot/log"

	"shop_bot/config"
	"shop_bot/storage"
)

type Bot struct {
	Config *config.Config

	Storage *storage.Storage

	TelegramAPI *tg.Client
	Sender      *message.Sender
}

func RunBot(cfg *config.Config) (err error) {
	b := &Bot{
		Config: cfg,
	}

	b.Storage, err = storage.NewStorage(b.Config.DB)
	if err != nil {
		log.Error("can't create new storage", err)
	}

	log.Info("starting order deactivation worker...")
	go b.orderDeactivationWorker()

	log.Info("creating telegram tools for bot...")
	dispatcher := tg.NewUpdateDispatcher()
	opts := telegram.Options{
		UpdateHandler: dispatcher,
	}

	log.Info("setting up env variables...")

	os.Setenv("BOT_TOKEN", b.Config.BotToken)
	os.Setenv("APP_ID", b.Config.ApiID)
	os.Setenv("APP_HASH", b.Config.ApiHash)

	log.Info("starting bot...")
	if err = telegram.BotFromEnvironment(context.Background(), opts, func(ctx context.Context, client *telegram.Client) error {
		b.TelegramAPI = tg.NewClient(client)
		b.Sender = message.NewSender(b.TelegramAPI)

		log.Info("setting up callback handlers...")

		dispatcher.OnBotCallbackQuery(func(ctx context.Context, e tg.Entities, update *tg.UpdateBotCallbackQuery) error {
			log.Info("callback query: update=%v data=%s", update, update.Data)

			callbackInfo := &callbackInfo{
				e:      e,
				update: update,
			}

			go func() {
				defer b.panicHandler()
				if err := b.callbackMapping(ctx, string(update.Data), callbackInfo); err != nil {
					log.Error("error occurred in callback", err)
				}
			}()

			return nil
		})

		log.Info("setting up new message handlers...")

		dispatcher.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
			log.Info("new message: update=%v", update)

			messageInfo := &messageInfo{
				e:      e,
				update: update,
			}

			go func() {
				defer b.panicHandler()
				if err := b.messageMapping(ctx, messageInfo); err != nil {
					log.Error("error occurred in message handler", err)
				}
			}()

			return nil
		})

		log.Info("finished setting up bot")

		return nil
	}, telegram.RunUntilCanceled); err != nil {
		log.Error("can't create bot from env", err)
		return nil
	}

	log.Info("bot started")

	return nil
}

func (b *Bot) panicHandler() {
	if err := recover(); err != nil {
		log.Info("recovered from panic: %v", err)
	}
}
