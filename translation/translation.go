package translation

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Translation interface {
	Trans(key string, args map[string]interface{}, languages ...string) string
	GetLocalization(lang string) *i18n.Localizer
}

var (
	trans *translation
	once  sync.Once
)

// Trans is a helper function that translates a message.
func Trans(key string, args map[string]interface{}, languages ...string) string {
	return trans.Trans(key, args, languages...)
}

type translation struct {
	config         Config
	acceptLanguage *i18n.Localizer
	bundle         *i18n.Bundle
}

// NewTranslation creates a new translation instance.
func NewTranslation(c Config) Translation {
	once.Do(func() {
		trans = &translation{
			config: c,
		}

		trans.createLocalePathIfNotExists(c.PathLocale)
		trans.bundle = i18n.NewBundle(language.English)
		trans.bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

		if err := trans.walkingInLocalePath(c.PathLocale); err != nil {
			log.Printf("Failed to walk in locale path: %s", err.Error())
		}
	})

	return trans
}

// GetLocalization initializes the localizer with the desired language.
func (t translation) GetLocalization(lang string) *i18n.Localizer {
	if lang == "" {
		lang = t.config.Locale
	}

	tag, err := language.Parse(lang)
	if err != nil {
		tag = language.English
		log.Printf("Failed to parse language tag: %s", err.Error())
	}

	t.acceptLanguage = i18n.NewLocalizer(t.bundle, tag.String())

	return t.acceptLanguage
}

// Trans is a helper function that translates a message.
func (t translation) Trans(key string, args map[string]interface{}, languages ...string) string {
	config := &i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: args,
	}

	if len(languages) > 0 {
		t.acceptLanguage = t.GetLocalization(languages[0])
	} else {
		t.acceptLanguage = t.GetLocalization(t.config.Locale)
	}

	message, err := t.acceptLanguage.Localize(config)
	if err != nil {
		defaultLang := i18n.NewLocalizer(t.bundle)
		if message, err = defaultLang.Localize(config); err != nil {
			log.Printf("Failed to localize message: %s", err.Error())
			return key
		} else {
			return message
		}
	}

	return message
}

// createLocaleDirectory creates a directory for the locale files.
func (t translation) createLocalePathIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Folder fa did not exist! We made it by default :)")
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			log.Printf("Failed to create locale path: %s", err.Error())
			return
		}
	}

	if _, err := os.Stat(path + "/fa.json"); !os.IsNotExist(err) {
		return
	}

	if _, err := os.Create(path + "/fa.json"); err != nil {
		log.Printf("Failed to create the default locale file: %s", err.Error())
		return
	} else {
		log.Printf("File fa.json did not exist! We made it by default :)")
	}
}

// walkingInLocalePath walks in the locale path and loads the message files.
func (t translation) walkingInLocalePath(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Failed to walk in locale path: %s", err.Error())
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			if _, err = t.bundle.LoadMessageFile(path); err != nil {
				log.Printf("Failed to walk in locale path: %s", err.Error())
			}
		}

		return nil
	})
}
