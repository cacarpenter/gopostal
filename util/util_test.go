package util

import (
	"fmt"
	"testing"
)

func TestSortArray(t *testing.T) {
	keys := []string{"ddd", "lll", "ppp", "www", "aaa"}
	unsorted := make([][]string, len(keys))
	for i, k := range keys {
		unsorted[i] = make([]string, 2)
		unsorted[i][0] = k
		unsorted[i][1] = fmt.Sprintf("val%d", i)
	}

	sorted := SortArray(unsorted)
	if len(sorted) != 5 {
		t.Fatal("sorted array length unexpected 5", len(sorted))
	}
	for i, kv := range sorted {
		if len(kv) != 2 {
			t.Fatalf("unexpected size %d %d\n", i, len(kv))
		}
	}
	if sorted[0][0] != "aaa" {
		t.Fatal("1st key unexpected")
	}
	if sorted[0][1] != "val4" {
		t.Fatal("1st val unexpected", sorted[0][1])
	}
	if sorted[1][0] != "ddd" {
		t.Fatal("2nd key unexpected")
	}
	if sorted[1][1] != "val0" {
		t.Fatal("2nd val unexpected")
	}
	if sorted[2][0] != "lll" {
		t.Fatal("3rd key unexpected")
	}
	if sorted[2][1] != "val1" {
		t.Fatal("3rd val unexpected")
	}
	if sorted[3][0] != "ppp" {
		t.Fatal("4th key unexpected")
	}
	if sorted[3][1] != "val2" {
		t.Fatal("4th val unexpected")
	}
	if sorted[4][0] != "www" {
		t.Fatal("5th key unexpected")
	}
	if sorted[4][1] != "val3" {
		t.Fatal("5th val unexpected")
	}
}
