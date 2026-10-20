package architecture_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
)

const module = "github.com/gabrzb/chess-go/"

// Each production package declares the project packages it may import.
var dependencies = map[string][]string{
	"cmd/chess-tui":                     {"internal/delivery/tui", "internal/domain", "internal/engine"},
	"internal/domain":                   {},
	"internal/domain/attacks":           {"internal/domain"},
	"internal/domain/movegen":           {"internal/domain"},
	"internal/domain/rules":             {"internal/domain", "internal/domain/attacks", "internal/domain/movegen"},
	"internal/domain/game":              {"internal/domain", "internal/domain/rules"},
	"internal/engine":                   {"internal/domain", "internal/domain/rules"},
	"internal/application":              {"internal/domain", "internal/domain/game"},
	"internal/delivery/tui":             {"internal/domain", "internal/domain/game", "internal/engine", "internal/delivery/tui/application", "internal/delivery/tui/match"},
	"internal/delivery/tui/application": {"internal/domain", "internal/domain/game", "internal/engine", "internal/delivery/tui/match", "internal/delivery/tui/render"},
	"internal/delivery/tui/match":       {"internal/application", "internal/domain", "internal/domain/game", "internal/delivery/tui/render"},
	"internal/delivery/tui/render":      {"internal/domain", "internal/domain/game"},
}

func allowedImport(pkg, imported string) bool {
	allowed, known := dependencies[pkg]
	if !known {
		return false
	}
	if strings.HasPrefix(imported, module) {
		return slices.Contains(allowed, strings.TrimPrefix(imported, module))
	}
	if pkg == "internal/application" || pkg == "internal/domain" || strings.HasPrefix(pkg, "internal/domain/") {
		return !strings.Contains(strings.Split(imported, "/")[0], ".")
	}
	return true
}

func TestProductionDependencies(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	files := 0
	for _, dir := range []string{"cmd", "internal"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			files++
			relative, err := filepath.Rel(root, filepath.Dir(path))
			if err != nil {
				return err
			}
			pkg := filepath.ToSlash(relative)
			if _, ok := dependencies[pkg]; !ok {
				t.Errorf("%s: declare dependencies for the new production package", path)
			}
			file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
			if err != nil {
				return err
			}
			for _, spec := range file.Imports {
				imported, err := strconv.Unquote(spec.Path.Value)
				if err != nil {
					return err
				}
				if !allowedImport(pkg, imported) {
					t.Errorf("%s: forbidden import %s", path, imported)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if files == 0 {
		t.Fatal("no production files checked")
	}
}

func TestDependencyPolicyRejectsBrokenBoundaries(t *testing.T) {
	for _, tc := range []struct{ pkg, imported string }{
		{"internal/domain", "charm.land/bubbletea/v2"},
		{"internal/application", "charm.land/bubbletea/v2"},
		{"internal/application", module + "internal/engine"},
		{"internal/domain", module + "internal/domain/rules"},
		{"internal/engine", module + "internal/delivery/tui"},
		{"internal/delivery/tui/match", module + "internal/engine"},
		{"internal/delivery/tui/render", module + "internal/domain/rules"},
		{"internal/domain/game", module + "internal/domain/tests/game"},
	} {
		if allowedImport(tc.pkg, tc.imported) {
			t.Errorf("accepted %s -> %s", tc.pkg, tc.imported)
		}
	}
}
