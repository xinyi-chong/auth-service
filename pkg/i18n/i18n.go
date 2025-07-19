package locale

import (
	"auth-service/internal/shared/consts"
	"auth-service/pkg/logger"
	"context"
	"embed"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

type TemplateData map[string]interface{}

//go:embed *.toml
var fs embed.FS

var (
	bundle *i18n.Bundle
)

func Init() error {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	if _, err := bundle.LoadMessageFileFS(fs, "active.en.toml"); err != nil {
		return err
	}
	if _, err := bundle.LoadMessageFileFS(fs, "active.zh.toml"); err != nil {
		return err
	}

	return nil
}

func GetLocalizer(langs ...string) *i18n.Localizer {
	return i18n.NewLocalizer(bundle, langs...)
}

func Translate(c context.Context, messageID string, templateData ...TemplateData) string {
	localizer, ok := c.Value(consts.Localizer).(*i18n.Localizer)
	if !ok || localizer == nil {
		logger.Log.Debug("Translate: no localizer found in context")
		return messageID
	}

	var data TemplateData
	if len(templateData) > 0 {
		data = templateData[0]
	}

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		logger.Log.Debug(
			"Translate: translation missing",
			zap.String("messageID", messageID),
			zap.Error(err),
		)
		return messageID // Fallback to message ID
	}
	return message
}
