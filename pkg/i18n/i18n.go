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
	"strings"
	"sync"
)

type TemplateData map[string]interface{}
type Category string
type Message string

const (
	CategorySuccess Category = "success"
	CategoryError   Category = "error"
)

//go:embed *.toml
var fs embed.FS

var (
	bundle     *i18n.Bundle
	bundleOnce sync.Once
)

func Init() error {
	var initErr error
	bundleOnce.Do(func() {
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

		entries, err := fs.ReadDir(".")
		if err != nil {
			initErr = err
			return
		}

		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".toml") {
				if _, err := bundle.LoadMessageFileFS(fs, entry.Name()); err != nil {
					initErr = err
					return
				}
			}
		}
	})
	return initErr
}

func GetLocalizer(langs ...string) *i18n.Localizer {
	return i18n.NewLocalizer(bundle, langs...)
}

func Translate(c context.Context, category Category, messageID string, templateData ...TemplateData) string {
	localizer, ok := c.Value(consts.Localizer).(*i18n.Localizer)
	if !ok || localizer == nil {
		logger.Error("Translate: no localizer found in context")
		return messageID
	}

	var data TemplateData
	if len(templateData) > 0 {
		data = templateData[0]
	}

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    string(category) + "." + messageID,
		TemplateData: data,
	})
	if err != nil {
		logger.Warn(
			"Translate: translation missing",
			zap.String("messageID", messageID),
			zap.Error(err),
		)
		return messageID
	}
	return message
}
