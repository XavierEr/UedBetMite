package uedBetDataJson

import "testing"

func TestExtractTotalPages(t *testing.T) {
	var floatTotalPages float64 = 3.2
	var expected int = 3

	if intTotalPages := extractTotalPages(floatTotalPages); intTotalPages != expected {
		t.Errorf("extractTotalPages(%v) = %v, expected %v", floatTotalPages, intTotalPages, expected)
	}

}
