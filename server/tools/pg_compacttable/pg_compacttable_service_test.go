package pg_compacttable

import (
	"ivory/plugins/database"
	"testing"
)

func TestSupportedFeatures(t *testing.T) {
	// NOTE: method under test does not touch dependencies, so the zero
	// service is enough and avoids the initializer goroutine in NewService
	s := &Service{}

	t.Run("postgres gets all features", func(t *testing.T) {
		features := s.SupportedFeatures(database.POSTGRES)
		for feature, supported := range features {
			if !supported {
				t.Errorf("expected feature %v to be supported for postgres", feature)
			}
		}
		if len(features) != 4 {
			t.Fatalf("expected 4 features for postgres, got %d", len(features))
		}
	})

	t.Run("etcd gets no supported features", func(t *testing.T) {
		features := s.SupportedFeatures(database.ETCD)
		for feature, supported := range features {
			if supported {
				t.Errorf("expected feature %v to not be supported for etcd", feature)
			}
		}
		if len(features) != 4 {
			t.Fatalf("expected 4 features for etcd, got %d", len(features))
		}
	})
}
