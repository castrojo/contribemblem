package bungie

import (
	"testing"
)

func TestManifestStructures(t *testing.T) {
	// Basic test to ensure structures compile
	_ = manifestResponse{}
	_ = emblemData{}
}
