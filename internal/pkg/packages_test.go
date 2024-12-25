package pkg

import (
	"sort"
	"testing"
)

func TestSorting(t *testing.T) {
	packages := PackageImports{
		{Path: "context"},
		{Path: "log/slog"},
		{Path: "os"},
		{Path: "time"},
		{Path: "github.com/joho/godotenv"},
	}
	sort.Sort(packages)

	if packages.String() != `	"context"\n	"log/slog"\n	"os"\n	"time"\n\n	"github.com/joho/godotenv"` {
		t.Errorf("unexpected imports: %s", packages.String())
	}
}
