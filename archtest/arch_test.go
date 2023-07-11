package archtest

import (
	"github.com/matthewmcnew/archtest"
	"testing"
)

const (
	domainModels      = "bitcoin-exchange-rate/internal/model/..."
	presentationLayer = "bitcoin-exchange-rate/internal/handler/..."
	serviceLayer      = "bitcoin-exchange-rate/internal/service/..."
	persistenceLayer  = "bitcoin-exchange-rate/internal/repository/..."
	providerClient    = "bitcoin-exchange-rate/pkg/rate_providers/..."
	mailerClient      = "bitcoin-exchange-rate/pkg/mailer/..."
)

func TestPresentationLayerDependency(t *testing.T) {
	archtest.Package(t, presentationLayer).ShouldNotDependOn(
		providerClient,
		mailerClient,
	)
}

func TestServiceLayerDependency(t *testing.T) {
	archtest.Package(t, serviceLayer).ShouldNotDependOn(
		presentationLayer,
	)
}

func TestPersistenceLayerDependency(t *testing.T) {
	archtest.Package(t, persistenceLayer).ShouldNotDependOn(
		serviceLayer,
		providerClient,
		persistenceLayer,
		mailerClient,
	)
}

func TestProviderClientDependency(t *testing.T) {
	archtest.Package(t, providerClient).ShouldNotDependOn(
		presentationLayer,
		serviceLayer,
		mailerClient,
		persistenceLayer,
	)
}

func TestMailerClientDependency(t *testing.T) {
	archtest.Package(t, mailerClient).ShouldNotDependOn(
		presentationLayer,
		serviceLayer,
		providerClient,
		persistenceLayer,
		domainModels,
	)
}

func TestDomainModelsDependency(t *testing.T) {
	archtest.Package(t, domainModels).ShouldNotDependOn(
		presentationLayer,
		serviceLayer,
		providerClient,
		persistenceLayer,
		mailerClient,
	)
}
