package permission

import (
	"ivory/clients/storage"
	"ivory/core/config"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
)

func createTestPermissionService(t *testing.T) *Service {
	t.Helper()

	tmpDir, errDir := os.MkdirTemp("", "permission-service-test-*")
	if errDir != nil {
		t.Fatalf("failed to create temp dir: %v", errDir)
	}
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	db, errOpen := bolt.Open(filepath.Join(tmpDir, "test.db"), 0600, nil)
	if errOpen != nil {
		t.Fatalf("failed to open test database: %v", errOpen)
	}
	t.Cleanup(func() {
		db.Close()
	})

	return NewService(NewRepository(storage.NewDbBucket[PermissionMap](db, "Permission")))
}

func TestServiceCreateUserPermissions(t *testing.T) {
	s := createTestPermissionService(t)

	t.Run("the caller's default is what every feature starts at", func(t *testing.T) {
		perms, err := s.CreateUserPermissions("alice", NOT_PERMITTED)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, feature := range config.All {
			if perms[feature] != NOT_PERMITTED {
				t.Fatalf("expected NOT_PERMITTED for %s, got %v", feature, perms[feature])
			}
		}
	})

	// NOTE: this feature has never heard of a superuser - GRANTED arrives as an
	// argument from the user feature, which is the one that knows
	t.Run("a granted default hands the user everything", func(t *testing.T) {
		perms, err := s.CreateUserPermissions("root", GRANTED)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, feature := range config.All {
			if perms[feature] != GRANTED {
				t.Fatalf("expected GRANTED for %s, got %v", feature, perms[feature])
			}
		}
	})

	t.Run("creating twice keeps the answers already stored", func(t *testing.T) {
		if _, err := s.CreateUserPermissions("bob", NOT_PERMITTED); err != nil {
			t.Fatalf("failed to create bob: %v", err)
		}
		if err := s.ApproveUserPermissions("bob", []config.Feature{config.ViewClusterList}); err != nil {
			t.Fatalf("failed to approve: %v", err)
		}
		perms, err := s.CreateUserPermissions("bob", NOT_PERMITTED)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if perms[config.ViewClusterList] != GRANTED {
			t.Fatalf("expected the previously granted permission to be preserved, got %v", perms[config.ViewClusterList])
		}
	})

	t.Run("a feature the record does not name yet is filled in with the default", func(t *testing.T) {
		if _, err := s.CreateUserPermissions("carol", NOT_PERMITTED); err != nil {
			t.Fatalf("failed to create carol: %v", err)
		}
		stored, errGet := s.GetUserPermissions("carol", false, nil)
		if errGet != nil {
			t.Fatalf("failed to read carol: %v", errGet)
		}
		delete(stored, config.ViewClusterList)
		if err := s.UpdateUserPermissions("carol", stored); err != nil {
			t.Fatalf("failed to store a record missing a feature: %v", err)
		}

		perms, err := s.CreateUserPermissions("carol", GRANTED)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if perms[config.ViewClusterList] != GRANTED {
			t.Fatalf("expected the missing feature to be filled in with the default, got %v", perms[config.ViewClusterList])
		}
	})

	t.Run("empty username is rejected", func(t *testing.T) {
		if _, err := s.CreateUserPermissions("", NOT_PERMITTED); err != ErrUsernameCannotBeEmpty {
			t.Fatalf("expected ErrUsernameCannotBeEmpty, got %v", err)
		}
	})
}

// TestServiceCreateUserPermissionsRenamesStoredFeatures covers a permission map
// stored before view.node.platform became view.node.system: filling a record in
// must carry the grant over to the new key rather than dropping it and
// resetting the user to the default status.
func TestServiceCreateUserPermissionsRenamesStoredFeatures(t *testing.T) {
	s := createTestPermissionService(t)
	if _, err := s.CreateUserPermissions("alice", NOT_PERMITTED); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	stored, err := s.GetUserPermissions("alice", false, nil)
	if err != nil {
		t.Fatalf("failed to read seeded permissions: %v", err)
	}
	delete(stored, config.ViewNodeSystem)
	stored["view.node.platform"] = GRANTED
	if err := s.UpdateUserPermissions("alice", stored); err != nil {
		t.Fatalf("failed to store the legacy permission: %v", err)
	}

	if _, err := s.CreateUserPermissions("alice", NOT_PERMITTED); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	perms, errGet := s.GetUserPermissions("alice", false, nil)
	if errGet != nil {
		t.Fatalf("expected no error, got %v", errGet)
	}
	if perms[config.ViewNodeSystem] != GRANTED {
		t.Fatalf("expected the renamed feature to keep its grant, got %v", perms[config.ViewNodeSystem])
	}
	if _, ok := perms["view.node.platform"]; ok {
		t.Fatal("expected the old key to be gone after normalization")
	}
}

func TestServiceGetUserPermissions(t *testing.T) {
	s := createTestPermissionService(t)
	if _, err := s.CreateUserPermissions("alice", NOT_PERMITTED); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	t.Run("allowAll bypasses lookup and grants everything", func(t *testing.T) {
		perms, err := s.GetUserPermissions("", true, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, feature := range config.All {
			if perms[feature] != GRANTED {
				t.Fatalf("expected GRANTED for %s, got %v", feature, perms[feature])
			}
		}
	})

	t.Run("existing user is fetched", func(t *testing.T) {
		perms, err := s.GetUserPermissions("alice", false, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(perms) != len(config.All) {
			t.Fatalf("expected %d features, got %d", len(config.All), len(perms))
		}
	})

	t.Run("unknown user fails", func(t *testing.T) {
		if _, err := s.GetUserPermissions("unknown", false, nil); err != ErrUserPermissionsNotFound {
			t.Fatalf("expected ErrUserPermissionsNotFound, got %v", err)
		}
	})

	t.Run("empty username fails", func(t *testing.T) {
		if _, err := s.GetUserPermissions("", false, nil); err != ErrUsernameCannotBeEmpty {
			t.Fatalf("expected ErrUsernameCannotBeEmpty, got %v", err)
		}
	})
}

func TestServiceRequestApproveRejectUserPermissions(t *testing.T) {
	s := createTestPermissionService(t)
	if _, err := s.CreateUserPermissions("alice", NOT_PERMITTED); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	t.Run("request marks features pending", func(t *testing.T) {
		if err := s.RequestUserPermissions("alice", []config.Feature{config.ViewClusterList}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		perms, _ := s.GetUserPermissions("alice", false, nil)
		if perms[config.ViewClusterList] != PENDING {
			t.Fatalf("expected PENDING, got %v", perms[config.ViewClusterList])
		}
	})

	t.Run("approve grants the feature", func(t *testing.T) {
		if err := s.ApproveUserPermissions("alice", []config.Feature{config.ViewClusterList}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		perms, _ := s.GetUserPermissions("alice", false, nil)
		if perms[config.ViewClusterList] != GRANTED {
			t.Fatalf("expected GRANTED, got %v", perms[config.ViewClusterList])
		}
	})

	t.Run("reject denies the feature", func(t *testing.T) {
		if err := s.RejectUserPermissions("alice", []config.Feature{config.ViewClusterList}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		perms, _ := s.GetUserPermissions("alice", false, nil)
		if perms[config.ViewClusterList] != NOT_PERMITTED {
			t.Fatalf("expected NOT_PERMITTED, got %v", perms[config.ViewClusterList])
		}
	})

	t.Run("invalid feature is rejected and reported per-feature", func(t *testing.T) {
		err := s.ApproveUserPermissions("alice", []config.Feature{"not-a-real-feature"})
		if err == nil {
			t.Fatalf("expected an error for an invalid feature")
		}
	})

	t.Run("aggregates errors across multiple features but keeps applying valid ones", func(t *testing.T) {
		err := s.ApproveUserPermissions("alice", []config.Feature{config.ViewClusterList, "not-a-real-feature"})
		if err == nil {
			t.Fatalf("expected an aggregated error")
		}
		perms, _ := s.GetUserPermissions("alice", false, nil)
		if perms[config.ViewClusterList] != GRANTED {
			t.Fatalf("expected the valid feature to still be applied, got %v", perms[config.ViewClusterList])
		}
	})

	t.Run("empty username is rejected", func(t *testing.T) {
		if err := s.ApproveUserPermissions("", []config.Feature{config.ViewClusterList}); err == nil {
			t.Fatalf("expected an error for an empty username")
		}
	})
}

// TestServiceApproveRejectUnknownUsername covers a name arriving off the
// request path, where nothing guarantees a record exists behind it. A key is
// the username itself now, so anything at all can be asked for - and a name
// nobody holds has to be an error rather than a record quietly written.
func TestServiceApproveRejectUnknownUsername(t *testing.T) {
	s := createTestPermissionService(t)
	features := []config.Feature{config.ViewClusterList}

	t.Run("a username nobody holds is refused", func(t *testing.T) {
		if err := s.ApproveUserPermissions("nobody", features); err == nil {
			t.Errorf("ApproveUserPermissions(\"nobody\") returned no error")
		}
		if err := s.RejectUserPermissions("nobody", features); err == nil {
			t.Errorf("RejectUserPermissions(\"nobody\") returned no error")
		}
	})

	// NOTE: a colon used to separate the authority from the name, and a
	// distinguished name is full of them - it is an ordinary username now
	t.Run("a username containing a colon is an ordinary name", func(t *testing.T) {
		if _, err := s.CreateUserPermissions("cn=bob,ou=eng:1", NOT_PERMITTED); err != nil {
			t.Fatalf("failed to seed: %v", err)
		}
		if err := s.ApproveUserPermissions("cn=bob,ou=eng:1", features); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		perms, _ := s.GetUserPermissions("cn=bob,ou=eng:1", false, nil)
		if perms[config.ViewClusterList] != GRANTED {
			t.Errorf("expected GRANTED, got %v", perms[config.ViewClusterList])
		}
	})
}

func TestServiceDeleteUserPermissions(t *testing.T) {
	s := createTestPermissionService(t)
	if _, err := s.CreateUserPermissions("alice", NOT_PERMITTED); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	t.Run("empty username is rejected", func(t *testing.T) {
		if err := s.DeleteUserPermissions(""); err != ErrUsernameCannotBeEmpty {
			t.Fatalf("expected ErrUsernameCannotBeEmpty, got %v", err)
		}
	})

	t.Run("existing user is deleted", func(t *testing.T) {
		if err := s.DeleteUserPermissions("alice"); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, err := s.GetUserPermissions("alice", false, nil); err == nil {
			t.Fatalf("expected an error getting a deleted user")
		}
	})
}

func TestServiceUpdateUserPermissions(t *testing.T) {
	s := createTestPermissionService(t)

	t.Run("empty username is rejected", func(t *testing.T) {
		if err := s.UpdateUserPermissions("", PermissionMap{}); err != ErrUsernameCannotBeEmpty {
			t.Fatalf("expected ErrUsernameCannotBeEmpty, got %v", err)
		}
	})

	t.Run("valid username persists the map", func(t *testing.T) {
		perms := PermissionMap{config.ViewClusterList: GRANTED}
		if err := s.UpdateUserPermissions("alice", perms); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		got, err := s.GetUserPermissions("alice", false, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got[config.ViewClusterList] != GRANTED {
			t.Fatalf("expected GRANTED, got %v", got[config.ViewClusterList])
		}
	})
}

func TestServiceGetAllUserPermissions(t *testing.T) {
	s := createTestPermissionService(t)
	if _, err := s.CreateUserPermissions("bob", NOT_PERMITTED); err != nil {
		t.Fatalf("failed to seed bob: %v", err)
	}
	if _, err := s.CreateUserPermissions("alice", NOT_PERMITTED); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	all, err := s.GetAllUserPermissions()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 users, got %d", len(all))
	}
	if all[0].Username != "alice" || all[1].Username != "bob" {
		t.Fatalf("expected results sorted by username, got %v", all)
	}
}

func TestServiceDeleteAll(t *testing.T) {
	s := createTestPermissionService(t)
	if _, err := s.CreateUserPermissions("alice", NOT_PERMITTED); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}
	if err := s.DeleteAll(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	all, err := s.GetAllUserPermissions()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(all) != 0 {
		t.Fatalf("expected no users, got %v", all)
	}
}

func TestPermissionMapRenamed(t *testing.T) {
	tests := []struct {
		name     string
		stored   PermissionMap
		expected PermissionMap
	}{
		{
			name:     "a renamed key carries its status over",
			stored:   PermissionMap{"view.node.platform": GRANTED},
			expected: PermissionMap{config.ViewNodeSystem: GRANTED},
		},
		{
			name:     "a current key is left alone",
			stored:   PermissionMap{config.ViewClusterList: PENDING},
			expected: PermissionMap{config.ViewClusterList: PENDING},
		},
		{
			// NOTE: map iteration order is random, so the current key has to win
			// deterministically rather than by whichever is visited last
			name:     "the current key wins over the old one",
			stored:   PermissionMap{"view.node.platform": GRANTED, config.ViewNodeSystem: NOT_PERMITTED},
			expected: PermissionMap{config.ViewNodeSystem: NOT_PERMITTED},
		},
		{
			name:     "an unknown key is passed through untouched",
			stored:   PermissionMap{"view.node.nonsense": GRANTED},
			expected: PermissionMap{"view.node.nonsense": GRANTED},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.stored.renamed()
			if len(got) != len(tt.expected) {
				t.Fatalf("renamed() = %v, want %v", got, tt.expected)
			}
			for feature, status := range tt.expected {
				if got[feature] != status {
					t.Errorf("renamed()[%q] = %v, want %v", feature, got[feature], status)
				}
			}
		})
	}
}
