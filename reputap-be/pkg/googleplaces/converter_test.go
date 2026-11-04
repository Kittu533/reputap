package googleplaces

import (
	"context"
	"testing"
)

func TestHexToPlaceID(t *testing.T) {
	pid, err := HexToPlaceID("0x2e7af9ff3fe94787", "0x3e3e9169605e362c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "ChIJh0fpP__5ei4RLDZeYGmRPj4"
	if pid != expected {
		t.Errorf("got %s, expected %s", pid, expected)
	}
}

func TestToDirectReviewURL(t *testing.T) {
	ctx := context.Background()
	input := "https://maps.app.goo.gl/daua4USiB8ezUxfV7"
	got := ToDirectReviewURL(ctx, input)
	expected := "https://search.google.com/local/writereview?placeid=ChIJh0fpP__5ei4RLDZeYGmRPj4"
	if got != expected {
		t.Errorf("got %s, expected %s", got, expected)
	}
}
