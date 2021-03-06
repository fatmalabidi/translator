package main

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/patrickmn/go-cache"
	"golang.org/x/text/language"
	"time"
)

// TODO move to a config file
const (
	DEFAUL_CACHE_EXPIRARION = 5 * time.Minute
	CACHE_PTGE_EXPIRED      = 10 * time.Minute
	EXPIRE_CACHE_DATA       = 1
)

type TranslatorWithCache struct {
	translator       *randomTranslator
	cache            *cache.Cache
	expirationMethod time.Duration // defines when to purge the expired data in the cache
}

func NewTranslatorWithCache() *TranslatorWithCache {
	t := newRandomTranslator(
		100*time.Millisecond,
		500*time.Millisecond,
		0.1,
	)
		c := cache.New(DEFAUL_CACHE_EXPIRARION, CACHE_PTGE_EXPIRED)

	return &TranslatorWithCache{
		translator:       t,
		cache:            c,
		expirationMethod: time.Duration(EXPIRE_CACHE_DATA),
	}
}

// Translate returns translation string or error. It checks the cache first for the requested
func (t *TranslatorWithCache) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	// check if the data exists in the cache.
	// the cache key pattern is "language:data" example "en:test"  to translate the word "test" to "english"
	key := fmt.Sprintf("%s-%s", language.English, data)
	cachedData, exist := t.cache.Get(key)
	if !exist {
		retries := func() error {
			translatedData, err := t.translator.Translate(ctx, from, to, data)
			if err == nil {
				t.cache.Set(key, translatedData, t.expirationMethod)
				return nil
			}
			return err
		}
		err := backoff.Retry(retries, backoff.NewExponentialBackOff())
		if err != nil {
			return "", err
		}
		translatedData, _ := t.cache.Get(key)
		return translatedData.(string), nil
	}
	return fmt.Sprintf("%v", cachedData), nil
}
