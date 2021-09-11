package main

import (
	"context"
	"fmt"
	"github.com/patrickmn/go-cache"
	"golang.org/x/text/language"
	"time"
)

// TODO move to a config file
const (
	DEFAUL_CACHE_EXPIRARION = 5 * time.Minute
	CACHE_PTGE_EXPIRED      = 10 * time.Minute
 )

type TranslatorWithCache struct {
	translator *randomTranslator
	cache      *cache.Cache
}

func NewTranslatorWithCache() *TranslatorWithCache {
	t := newRandomTranslator(
		100*time.Millisecond,
		500*time.Millisecond,
		0.1,
	)
	c := cache.New(DEFAUL_CACHE_EXPIRARION, CACHE_PTGE_EXPIRED)
	return &TranslatorWithCache{
		translator: t,
		cache:      c,
	}
}

// Translate returns translation string or error. It checks the cache first for the requested
func (t *TranslatorWithCache) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	// check if the data exists in the cache.
	// the cache key pattern is "language:data" example "en:test"  to translate the word "test" to "english"
	key := fmt.Sprintf("%s-%s", language.English, data)
	cachedData, exist := t.cache.Get(key)
	if !exist {
		translatedData, err := t.translator.Translate(ctx, from, to, data)
		if err == nil {
			t.cache.Set(key, translatedData, cache.NoExpiration)
			return translatedData, nil
		}
		return "", err
	}
	return fmt.Sprintf("%v", cachedData), nil
}
