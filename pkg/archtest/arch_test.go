package archtest

import (
	"github.com/matthewmcnew/archtest"
	"testing"
)

const (
	notificationModule = "bitcoin-exchange-rate/modules/notification_module/..."
	rateModule         = "bitcoin-exchange-rate/modules/rate_module/..."
	pkg                = "bitcoin-exchange-rate/pkg"
)

func TestRateModuleDependency(t *testing.T) {
	archtest.Package(t, rateModule).ShouldNotDependOn(
		notificationModule,
		pkg,
	)
}

func TestNotificationModuleDependency(t *testing.T) {
	archtest.Package(t, notificationModule).ShouldNotDependOn(
		rateModule,
		pkg,
	)
}
