package main

import (
	"context"
	"fmt"
	"golang.org/x/text/language"
	"testing"
)

func TestTranslate(t *testing.T) {
	t.Run("test translator", func(t *testing.T) {
		s := NewService()
		data := "test"
		ctx := context.Background()
		_, err := s.translator.Translate(ctx, language.English, language.Japanese, data)
		if err != nil {
			t.Error("expected success got error", err)
		}
		key := fmt.Sprintf("%s-%s", language.English, data)
		_, exist := s.translator.cache.Get(key)
		if !exist {
			t.Error("expected data to be cached got nothing")
		}
	})
}
