package utils

import "testing"

func TestSumInts(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want int
	}{
		{"empty", []int{}, 0},
		{"one", []int{5}, 5},
		{"many", []int{1, 2, 3, 4}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SumInts(tt.in)
			if got != tt.want {
				t.Fatalf("SumInts(%v) = %d; want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestMapKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	ks := MapKeys(m)
	if len(ks) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(ks))
	}
}

func TestErrorExample(t *testing.T) {
	_, err := ErrorExample("   ")
	if err == nil {
		t.Fatalf("expected error for empty input")
	}
	out, err := ErrorExample("go")
	if err != nil || out != "GO" {
		t.Fatalf("expected GO, got %q, err=%v", out, err)
	}
}
