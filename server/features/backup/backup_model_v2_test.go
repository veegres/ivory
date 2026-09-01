package backup

import "testing"

// TestBackupV2ShapeIsFrozen holds the sacred rule for V2. It stays editable
// only while V2 has shipped in prereleases alone - a stable major or minor tag
// freezes it, and from then on a shape change is a BackupV3, never an edit
// here. Regenerating testdata/ivory.v2.bak to make this pass is the one thing
// that defeats the test: after the freeze, the file is the specification.
func TestBackupV2ShapeIsFrozen(t *testing.T) {
	assertShapeIsFrozen(t, "ivory.v2.bak", &BackupV2{})
}
