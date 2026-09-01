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

func TestServiceSetSuperusers(t *testing.T) {
	s := createTestPermissionService(t)

	t.Run("empty list is rejected", func(t *testing.T) {
		if err := s.SetSuperusers(nil); err != ErrAtLeastOneSuperuser {
			t.Fatalf("expected ErrAtLeastOneSuperuser, got %v", err)
		}
	})

	t.Run("empty username in the list is rejected", func(t *testing.T) {
		if err := s.SetSuperusers([]string{"admin", ""}); err != ErrSuperusersCannotHaveEmptyName {
			t.Fatalf("expected ErrSuperusersCannotHaveEmptyName, got %v", err)
		}
	})

	t.Run("valid list is accepted and normalizes the database", func(t *testing.T) {
		if _, err := s.CreateUserPermissions("basic", "alice"); err != nil {
			t.Fatalf("failed to seed alice: %v", err)
		}
		if err := s.SetSuperusers([]string{"admin"}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestServiceDeleteAdmins(t *testing.T) {
	s := createTestPermissionService(t)
	if err := s.SetSuperusers([]string{"admin"}); err != nil {
		t.Fatalf("failed to set superusers: %v", err)
	}
	s.DeleteAdmins()
	if len(s.superusers) != 0 {
		t.Fatalf("expected superusers to be cleared, got %v", s.superusers)
	}
}

func TestServiceCreateUserPermissions(t *testing.T) {
	s := createTestPermissionService(t)

	t.Run("regular user gets NOT_PERMITTED for every feature", func(t *testing.T) {
		perms, err := s.CreateUserPermissions("basic", "alice")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, feature := range config.All {
			if perms[feature] != NOT_PERMITTED {
				t.Fatalf("expected NOT_PERMITTED for %s, got %v", feature, perms[feature])
			}
		}
	})

	t.Run("superuser gets GRANTED for every feature", func(t *testing.T) {
		if err := s.SetSuperusers([]string{"root"}); err != nil {
			t.Fatalf("failed to set superusers: %v", err)
		}
		perms, err := s.CreateUserPermissions("basic", "root")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		for _, feature := range config.All {
			if perms[feature] != GRANTED {
				t.Fatalf("expected GRANTED for %s, got %v", feature, perms[feature])
			}
		}
	})

	t.Run("creating twice returns the existing permissions instead of resetting them", func(t *testing.T) {
		if _, err := s.CreateUserPermissions("basic", "bob"); err != nil {
			t.Fatalf("failed to create bob: %v", err)
		}
		if err := s.ApproveUserPermissions("basic:bob", []config.Feature{config.ViewClusterList}); err != nil {
			t.Fatalf("failed to approve: %v", err)
		}
		perms, err := s.CreateUserPermissions("basic", "bob")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if perms[config.ViewClusterList] != GRANTED {
			t.Fatalf("expected the previously granted permission to be preserved, got %v", perms[config.ViewClusterList])
		}
	})

	t.Run("empty username is rejected", func(t *testing.T) {
		if _, err := s.CreateUserPermissions("basic", ""); err != ErrUsernameCannotBeEmpty {
			t.Fatalf("expected ErrUsernameCannotBeEmpty, got %v", err)
		}
	})

	t.Run("empty prefix is rejected", func(t *testing.T) {
		if _, err := s.CreateUserPermissions("", "alice"); err != ErrPrefixCannotBeEmpty {
			t.Fatalf("expected ErrPrefixCannotBeEmpty, got %v", err)
		}
	})
}

func TestServiceGetUserPermissions(t *testing.T) {
	s := createTestPermissionService(t)
	if _, err := s.CreateUserPermissions("basic", "alice"); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	t.Run("allowAll bypasses lookup and grants everything", func(t *testing.T) {
		perms, err := s.GetUserPermissions("", "", true)
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
		perms, err := s.GetUserPermissions("basic", "alice", false)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(perms) != len(config.All) {
			t.Fatalf("expected %d features, got %d", len(config.All), len(perms))
		}
	})

	t.Run("unknown user fails", func(t *testing.T) {
		if _, err := s.GetUserPermissions("basic", "unknown", false); err == nil {
			t.Fatalf("expected an error for an unknown user")
		}
	})
}

func TestServiceRequestApproveRejectUserPermissions(t *testing.T) {
	s := createTestPermissionService(t)
	if _, err := s.CreateUserPermissions("basic", "alice"); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	t.Run("request marks features pending", func(t *testing.T) {
		if err := s.RequestUserPermissions("basic", "alice", []config.Feature{config.ViewClusterList}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		perms, _ := s.GetUserPermissions("basic", "alice", false)
		if perms[config.ViewClusterList] != PENDING {
			t.Fatalf("expected PENDING, got %v", perms[config.ViewClusterList])
		}
	})

	t.Run("approve grants the feature", func(t *testing.T) {
		if err := s.ApproveUserPermissions("basic:alice", []config.Feature{config.ViewClusterList}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		perms, _ := s.GetUserPermissions("basic", "alice", false)
		if perms[config.ViewClusterList] != GRANTED {
			t.Fatalf("expected GRANTED, got %v", perms[config.ViewClusterList])
		}
	})

	t.Run("reject denies the feature", func(t *testing.T) {
		if err := s.RejectUserPermissions("basic:alice", []config.Feature{config.ViewClusterList}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		perms, _ := s.GetUserPermissions("basic", "alice", false)
		if perms[config.ViewClusterList] != NOT_PERMITTED {
			t.Fatalf("expected NOT_PERMITTED, got %v", perms[config.ViewClusterList])
		}
	})

	t.Run("cannot change permissions for superusers", func(t *testing.T) {
		if err := s.SetSuperusers([]string{"root"}); err != nil {
			t.Fatalf("failed to set superusers: %v", err)
		}
		if _, err := s.CreateUserPermissions("basic", "root"); err != nil {
			t.Fatalf("failed to seed root: %v", err)
		}
		err := s.ApproveUserPermissions("basic:root", []config.Feature{config.ViewClusterList})
		if err == nil {
			t.Fatalf("expected an error for changing a superuser's permissions")
		}
	})

	t.Run("invalid feature is rejected and reported per-feature", func(t *testing.T) {
		err := s.ApproveUserPermissions("basic:alice", []config.Feature{"not-a-real-feature"})
		if err == nil {
			t.Fatalf("expected an error for an invalid feature")
		}
	})

	t.Run("aggregates errors across multiple features but keeps applying valid ones", func(t *testing.T) {
		err := s.ApproveUserPermissions("basic:alice", []config.Feature{config.ViewClusterList, "not-a-real-feature"})
		if err == nil {
			t.Fatalf("expected an aggregated error")
		}
		perms, _ := s.GetUserPermissions("basic", "alice", false)
		if perms[config.ViewClusterList] != GRANTED {
			t.Fatalf("expected the valid feature to still be applied, got %v", perms[config.ViewClusterList])
		}
	})

	t.Run("empty permUsername is rejected", func(t *testing.T) {
		if err := s.ApproveUserPermissions("", []config.Feature{config.ViewClusterList}); err == nil {
			t.Fatalf("expected an error for an empty username")
		}
	})
}

// TestServiceApproveRejectRequireAPrefixedUsername covers a name arriving
// without its prefix. It comes straight off the request path, so nothing
// guarantees the shape getFullUsername writes - and the old code indexed the
// split blindly, panicking into an empty 500 rather than saying what was wrong.
func TestServiceApproveRejectRequireAPrefixedUsername(t *testing.T) {
	s := createTestPermissionService(t)
	features := []config.Feature{config.ViewClusterList}

	malformed := []struct {
		name     string
		username string
	}{
		{name: "no prefix at all", username: "alice"},
		{name: "empty prefix", username: ":alice"},
		{name: "empty username", username: "basic:"},
	}

	for _, tt := range malformed {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.ApproveUserPermissions(tt.username, features); err == nil {
				t.Errorf("ApproveUserPermissions(%q) returned no error", tt.username)
			}
			if err := s.RejectUserPermissions(tt.username, features); err == nil {
				t.Errorf("RejectUserPermissions(%q) returned no error", tt.username)
			}
		})
	}

	// NOTE: a colon in the username itself is only the separator once - the
	// prefix ends at the first one, and the rest is the name
	t.Run("a username containing a colon keeps its whole name", func(t *testing.T) {
		if _, err := s.CreateUserPermissions("ldap", "cn=bob,ou=eng"); err != nil {
			t.Fatalf("failed to seed: %v", err)
		}
		if err := s.ApproveUserPermissions("ldap:cn=bob,ou=eng", features); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		perms, _ := s.GetUserPermissions("ldap", "cn=bob,ou=eng", false)
		if perms[config.ViewClusterList] != GRANTED {
			t.Errorf("expected GRANTED, got %v", perms[config.ViewClusterList])
		}
	})
}

func TestServiceDeleteUserPermissions(t *testing.T) {
	s := createTestPermissionService(t)
	if _, err := s.CreateUserPermissions("basic", "alice"); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	t.Run("empty username is rejected", func(t *testing.T) {
		if err := s.DeleteUserPermissions(""); err != ErrUsernameCannotBeEmpty {
			t.Fatalf("expected ErrUsernameCannotBeEmpty, got %v", err)
		}
	})

	t.Run("existing user is deleted", func(t *testing.T) {
		if err := s.DeleteUserPermissions("basic:alice"); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, err := s.GetUserPermissions("basic", "alice", false); err == nil {
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
		if err := s.UpdateUserPermissions("basic:alice", perms); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		got, err := s.GetUserPermissions("basic", "alice", false)
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
	if _, err := s.CreateUserPermissions("basic", "bob"); err != nil {
		t.Fatalf("failed to seed bob: %v", err)
	}
	if _, err := s.CreateUserPermissions("basic", "alice"); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	all, err := s.GetAllUserPermissions()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 users, got %d", len(all))
	}
	if all[0].Username != "basic:alice" || all[1].Username != "basic:bob" {
		t.Fatalf("expected results sorted by username, got %v", all)
	}
}

func TestServiceDeleteAll(t *testing.T) {
	s := createTestPermissionService(t)
	if _, err := s.CreateUserPermissions("basic", "alice"); err != nil {
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

func TestServiceSetSuperusersNormalizesExistingUsers(t *testing.T) {
	s := createTestPermissionService(t)
	if _, err := s.CreateUserPermissions("basic", "alice"); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}
	if err := s.ApproveUserPermissions("basic:alice", []config.Feature{config.ViewClusterList}); err != nil {
		t.Fatalf("failed to approve: %v", err)
	}

	if err := s.SetSuperusers([]string{"admin"}); err != nil {
		t.Fatalf("failed to set superusers: %v", err)
	}

	perms, err := s.GetUserPermissions("basic", "alice", false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if perms[config.ViewClusterList] != GRANTED {
		t.Fatalf("expected the previously granted feature to be preserved, got %v", perms[config.ViewClusterList])
	}
}

// TestServiceSetSuperusersRenamesStoredFeatures covers a permission map stored
// before view.node.platform became view.node.system: normalizing must carry the
// grant over to the new key rather than dropping it and resetting the user to
// the default status.
func TestServiceSetSuperusersRenamesStoredFeatures(t *testing.T) {
	s := createTestPermissionService(t)
	if _, err := s.CreateUserPermissions("basic", "alice"); err != nil {
		t.Fatalf("failed to seed alice: %v", err)
	}

	stored, err := s.GetUserPermissions("basic", "alice", false)
	if err != nil {
		t.Fatalf("failed to read seeded permissions: %v", err)
	}
	delete(stored, config.ViewNodeSystem)
	stored["view.node.platform"] = GRANTED
	if err := s.UpdateUserPermissions("basic:alice", stored); err != nil {
		t.Fatalf("failed to store the legacy permission: %v", err)
	}

	if err := s.SetSuperusers([]string{"admin"}); err != nil {
		t.Fatalf("failed to set superusers: %v", err)
	}

	perms, err := s.GetUserPermissions("basic", "alice", false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if perms[config.ViewNodeSystem] != GRANTED {
		t.Fatalf("expected the renamed feature to keep its grant, got %v", perms[config.ViewNodeSystem])
	}
	if _, ok := perms["view.node.platform"]; ok {
		t.Fatal("expected the old key to be gone after normalization")
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
