package courier

import (
	"net/url"
	"testing"
	"time"
)

func TestParseQueryRejectsInvalidLevel(t *testing.T) {
	_, errorsByField := ParseQuery(url.Values{"level": {"2,9"}})
	if len(errorsByField["level"]) == 0 {
		t.Fatal("expected invalid level to be rejected")
	}
}

func TestParseQueryParsesMultipleLevels(t *testing.T) {
	params, errorsByField := ParseQuery(url.Values{"level": {"2,3"}})
	if len(errorsByField) > 0 {
		t.Fatalf("expected valid query, got %#v", errorsByField)
	}
	if len(params.Levels) != 2 || params.Levels[0] != 2 || params.Levels[1] != 3 {
		t.Fatalf("expected levels 2 and 3, got %#v", params.Levels)
	}
}

func TestParseQueryRejectsArbitrarySort(t *testing.T) {
	_, errorsByField := ParseQuery(url.Values{"sort": {"deleted_at"}})
	if len(errorsByField["sort"]) == 0 {
		t.Fatal("expected invalid sort to be rejected")
	}
}

func TestValidatorRejectsInvalidPayload(t *testing.T) {
	validationErrors := NewValidator().ValidatePayload(CourierPayload{Name: "", Level: 9})
	if len(validationErrors["name"]) == 0 || len(validationErrors["level"]) == 0 {
		t.Fatalf("expected name and level validation errors, got %#v", validationErrors)
	}
}

func TestParseRegisteredAtAcceptsDate(t *testing.T) {
	value := "2026-05-17"
	registeredAt, err := parseRegisteredAt(&value, mustDate(t, "2026-01-01"))
	if err != nil {
		t.Fatal(err)
	}
	if registeredAt.Format("2006-01-02") != "2026-05-17" {
		t.Fatalf("unexpected date %s", registeredAt)
	}
}

func mustDate(t *testing.T, raw string) time.Time {
	t.Helper()
	parsed, err := time.Parse("2006-01-02", raw)
	if err != nil {
		t.Fatal(err)
	}
	return parsed
}
