package pkg

import (
	"errors"
)

type ProviderNode interface {
	RateProvider
	SetNext(provider ProviderNode)
}

type ProvidersChain interface {
	RegisterProvider(name string, exchanger, next ProviderNode) error
	GetProvider(name string) ProviderNode
}

type providersChain struct {
	providers map[string]ProviderNode
}

func NewProvidersChain() ProvidersChain {
	return &providersChain{
		providers: make(map[string]ProviderNode),
	}
}

func (c *providersChain) RegisterProvider(name string, provider, next ProviderNode) error {
	if len(name) < 1 || provider == nil {
		return errors.New("invalid provider on input")
	}

	c.providers[name] = provider
	c.providers[name].SetNext(next)

	return nil
}

func (c *providersChain) GetProvider(name string) ProviderNode {
	return c.providers[name]
}
