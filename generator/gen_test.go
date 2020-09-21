package generator

import (
	"sort"
	"testing"
)

func TestListYamlFiles(t *testing.T) {
	expected := []string{
		"../test/specs/constraints.yaml",
		"../test/specs/dependsOn.yaml",
		"../test/specs/my.views.yaml",
		"../test/specs/noValidator.yaml",
		"../test/specs/requiredField.yaml",
		"../test/views/some_view.yml",
	}

	files, err := listYamlFiles("../test")
	if err != nil {
		t.Fatalf("should not error: %s", err)
	}

	sort.Strings(files)
	if len(files) != len(expected) {
		t.Fatalf("expected number of files %d found %d", len(expected), len(files))
	}

	for i, v := range files {
		if v != expected[i] {
			t.Fatalf("expected %s found %s", expected[i], v)
			break
		}
	}
}
