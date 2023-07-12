package pkg

import (
	"errors"
)

type ProvidersChain interface {
	RegisterProvider(name string, exchanger, next RateProviderNode) error
	GetProvider(name string) RateProviderNode
}

type providersChain struct {
	providers map[string]RateProviderNode
}

func NewProvidersChain() ProvidersChain {
	return &providersChain{
		providers: make(map[string]RateProviderNode),
	}
}

func (c *providersChain) RegisterProvider(name string, provider, next RateProviderNode) error {
	if len(name) < 1 || provider == nil {
		return errors.New("invalid provider on input")
	}

	c.providers[name] = provider
	c.providers[name].SetNext(next)

	return nil
}

func (c *providersChain) GetProvider(name string) RateProviderNode {
	return c.providers[name]
}
