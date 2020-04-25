package utils

import "testing"

func TestCreateHeader(t *testing.T) {
	t.Run("one to one same name", func(t *testing.T) {
		columns := []string{"a", "b", "c"}
		fieldColumnMap := map[string]string{"a": "a", "b": "b", "c": "c"}
		expected := map[string]int{"a": 0, "b": 1, "c": 2}

		found, err := CreateHeader(columns, fieldColumnMap)
		if err != nil {
			t.Error("must not returns error")
		}
		if len(expected) != len(found) {
			t.Errorf("different length: expected=%d found=%d",
				len(expected), len(found))
		}
		for k, v := range expected {
			if found[k] != v {
				t.Errorf("different value for key '%s': expected=%d found=%d",
					k, v, found[k])
			}
		}
	})

	t.Run("one to one different name", func(t *testing.T) {
		columns := []string{"a", "b", "c"}
		fieldColumnMap := map[string]string{"a": "b", "b": "c", "c": "a"}
		expected := map[string]int{"a": 1, "b": 2, "c": 0}

		found, err := CreateHeader(columns, fieldColumnMap)
		if err != nil {
			t.Error("must not returns error")
		}
		if len(expected) != len(found) {
			t.Errorf("different length: expected=%d found=%d",
				len(expected), len(found))
		}
		for k, v := range expected {
			if found[k] != v {
				t.Errorf("different value for key '%s': expected=%d found=%d",
					k, v, found[k])
			}
		}
	})

	t.Run("multiple fields to one column", func(t *testing.T) {
		columns := []string{"a", "b"}
		fieldColumnMap := map[string]string{"a": "a", "b": "a", "c": "b"}
		expected := map[string]int{"a": 0, "b": 0, "c": 1}

		found, err := CreateHeader(columns, fieldColumnMap)
		if err != nil {
			t.Error("must not returns error")
		}
		if len(expected) != len(found) {
			t.Errorf("different length: expected=%d found=%d",
				len(expected), len(found))
		}
		for k, v := range expected {
			if found[k] != v {
				t.Errorf("different value for key '%s': expected=%d found=%d",
					k, v, found[k])
			}
		}
	})

	t.Run("error when referencing to not exists column", func(t *testing.T) {
		columns := []string{"a", "b"}
		fieldColumnMap := map[string]string{"a": "a", "b": "b", "x": "c"}
		expected := "column 'c' for field 'x' not found"

		_, err := CreateHeader(columns, fieldColumnMap)
		if err == nil {
			t.Error("must returns error")
		}

		if expected != err.Error() {
			t.Errorf(`different error: expected="%s" found="%s"`,
				expected, err.Error())
		}
	})
}
