package deps

import (
	"testing"
)

func TestCsv2Postgres(t *testing.T) {
	dg := NewDepsGraph([]string{
		"noValidator",
		"constraints",
		"dependsOn",
	})

	err := dg.DependsOn("dependsOn", "noValidator")
	if err != nil {
		t.Errorf("should not error: %s", err)
	}

	err = dg.DependsOn("dependsOn", "constraints")
	if err != nil {
		t.Errorf("should not error: %s", err)
	}

	err = dg.Finalize()
	if err != nil {
		t.Errorf("should not error: %s", err)
	}

	p, err := dg.CreateOrder("dependsOn")
	expected := []string{"constraints", "noValidator", "dependsOn"}
	for i, exp := range expected {
		if exp != p[i] {
			t.Errorf("different result: expected=%+v found=%+v", expected, p)
		}
	}
}
