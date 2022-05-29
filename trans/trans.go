package trans

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"
)

type Translator interface {
	Translate(japanese string) (string, error)
}

type translatorImpl struct {
	translator *translate.Translate
}

func New() Translator {
	sess := session.Must(session.NewSession())
	translator := translate.New(sess)
	return &translatorImpl{translator}
}

func (t *translatorImpl) Translate(japanese string) (string, error) {
	source := "ja"
	target := "en"
	result, err := t.translator.Text(&translate.TextInput{
		SourceLanguageCode: &source,
		TargetLanguageCode: &target,
		Text:               &japanese,
	})
	if err != nil {
		return "", fmt.Errorf("failed to translate text; %s; %w", japanese, err)
	}
	return *result.TranslatedText, nil
}
