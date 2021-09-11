package main

// Service is a Translator user.
type Service struct {
	translator *TranslatorWithCache
}

func NewService() *Service {
	t := NewTranslatorWithCache()
	return &Service{
		translator: t,
	}
}
