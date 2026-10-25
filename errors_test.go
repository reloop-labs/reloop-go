package reloop

import (
	"encoding/json"
	"testing"
)

func TestAPIErrorBodyReadsFix(t *testing.T) {
	var body APIErrorBody
	if err := json.Unmarshal([]byte(`{"message":"Domain not verified","why":"No DKIM","fix":"Verify the domain","link":"https://reloop.sh/docs"}`), &body); err != nil {
		t.Fatal(err)
	}
	if body.Fix != "Verify the domain" || body.Tip != "Verify the domain" {
		t.Fatalf("body = %#v", body)
	}
}

func TestAPIErrorBodyAcceptsLegacyTip(t *testing.T) {
	var body APIErrorBody
	if err := json.Unmarshal([]byte(`{"message":"x","tip":"Try again"}`), &body); err != nil {
		t.Fatal(err)
	}
	if body.Fix != "Try again" || body.Tip != "Try again" {
		t.Fatalf("body = %#v", body)
	}
}
