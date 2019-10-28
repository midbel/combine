package combine

import (
	"fmt"
	"strings"
	"testing"
)

func TestCombineAndLink(t *testing.T) {
	var (
		abc  = Single([]string{"A", "B", "C"})
		one  = Single([]string{"1", "2", "3"})
		xyz  = Single([]string{"X", "Y"})
		two  = Single([]string{"11", "22"})
		data = [][]string{
			{"A", "1", "X", "11"},
			{"A", "1", "Y", "22"},
			{"B", "2", "X", "11"},
			{"B", "2", "Y", "22"},
			{"C", "3", "X", "11"},
			{"C", "3", "Y", "22"},
		}
	)
	src := CombineSources(LinkSources(abc, one), LinkSources(xyz, two))
	if err := testSources(src, data); err != nil {
		t.Errorf("combination failure: %s", err)
	}
}

func TestLinkStrings(t *testing.T) {
	data := []struct {
		Left  []string
		Right []string
		Data  [][]string
	}{
		{
			Left:  []string{"1", "2"},
			Right: []string{"A", "B", "C"},
			Data: [][]string{
				{"1", "A"},
				{"2", "B"},
			},
		},
	}
	for i, d := range data {
		src := LinkStrings(d.Left, d.Right)
		if err := testSources(src, d.Data); err != nil {
			t.Errorf("%d) linkage failure (%s, %s): %s", i+1, d.Left, d.Right, err)
		}
	}
}

func TestCombineStrings(t *testing.T) {
	data := []struct {
		Left  []string
		Right []string
		Data  [][]string
	}{
		{
			Left:  []string{"1", "2", "3"},
			Right: []string{"A", "B"},
			Data: [][]string{
				{"1", "A"},
				{"1", "B"},
				{"2", "A"},
				{"2", "B"},
				{"3", "A"},
				{"3", "B"},
			},
		},
	}
	for i, d := range data {
		src := CombineStrings(d.Left, d.Right)
		if err := testSources(src, d.Data); err != nil {
			t.Errorf("%d) combination failure (%s, %s): %s", i+1, d.Left, d.Right, err)
		}
	}
}

func testSources(src Source, data [][]string) error {
	for i := 0; ; i++ {
		switch set, err := src.Next(); err {
		case nil:
			if i >= len(data) {
				return fmt.Errorf("too many combinations generated (%d >= %d)", i, len(data))
			}
			got, want := strings.Join(set, "/"), strings.Join(data[i], "/")
			if got != want {
				return fmt.Errorf("mismatch: %s != %s", got, want)
			}
		case ErrDone:
			return nil
		default:
			return err
		}
	}
}
