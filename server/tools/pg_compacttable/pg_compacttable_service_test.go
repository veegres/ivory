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
		if len(features) != 4 {
			t.Fatalf("expected 4 features for postgres, got %d", len(features))
		}
	})

	t.Run("etcd gets no features", func(t *testing.T) {
		features := s.SupportedFeatures(database.ETCD)
		if len(features) != 0 {
			t.Fatalf("expected no features for etcd, got %d", len(features))
		}
	})
}
