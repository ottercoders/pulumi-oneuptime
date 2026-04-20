package resources

import (
	"testing"
)

func TestAttachRotation_Envelope(t *testing.T) {
	t.Parallel()
	data := map[string]interface{}{}
	attachRotation(data, "rotation", &Recurring{IntervalType: "Day", IntervalCount: 7})

	outer, ok := data["rotation"].(map[string]interface{})
	if !ok || outer["_type"] != "Recurring" {
		t.Fatalf("rotation envelope wrong: %v", data["rotation"])
	}
	value := outer["value"].(map[string]interface{})
	if value["intervalType"] != "Day" {
		t.Errorf("intervalType = %v", value["intervalType"])
	}
	countEnv, ok := value["intervalCount"].(map[string]interface{})
	if !ok || countEnv["_type"] != "PositiveNumber" || countEnv["value"] != 7 {
		t.Errorf("intervalCount envelope wrong: %v", value["intervalCount"])
	}
}

func TestAttachRotation_Nil(t *testing.T) {
	t.Parallel()
	data := map[string]interface{}{"keep": "me"}
	attachRotation(data, "rotation", nil)
	if _, present := data["rotation"]; present {
		t.Error("nil rotation should leave key absent")
	}
}

func TestAttachRestrictionTimes_WeeklyWithTypoKey(t *testing.T) {
	t.Parallel()
	data := map[string]interface{}{}
	attachRestrictionTimes(data, "restrictionTimes", &RestrictionTimes{
		RestrictionType: "Weekly",
		WeeklyRestrictionTimes: []WeeklyRestriction{{
			StartDay: "Monday", EndDay: "Friday",
			StartTime: "09:00", EndTime: "17:00",
		}},
	})

	outer := data["restrictionTimes"].(map[string]interface{})
	if outer["_type"] != "RestrictionTimes" {
		t.Errorf("_type = %v", outer["_type"])
	}
	value := outer["value"].(map[string]interface{})
	// Upstream typo preserved on wire.
	if value["restictionType"] != "Weekly" {
		t.Errorf("expected misspelled 'restictionType' key with value 'Weekly'; got %v", value)
	}
	if _, wrong := value["restrictionType"]; wrong {
		t.Error("should NOT emit corrected-spelling 'restrictionType' key on wire")
	}
	weekly := value["weeklyRestrictionTimes"].([]map[string]interface{})
	if len(weekly) != 1 {
		t.Fatalf("expected 1 weekly entry, got %d", len(weekly))
	}
	w := weekly[0]
	if w["startDay"] != "Monday" || w["endDay"] != "Friday" {
		t.Errorf("weekly entry wrong: %v", w)
	}
	// Day restriction not supplied → null on wire (RestrictionTimes.ts keeps
	// the key always present).
	if value["dayRestrictionTimes"] != nil {
		t.Errorf("dayRestrictionTimes should be null, got %v", value["dayRestrictionTimes"])
	}
}

func TestAttachRestrictionTimes_DailyOverrides(t *testing.T) {
	t.Parallel()
	data := map[string]interface{}{}
	attachRestrictionTimes(data, "restrictionTimes", &RestrictionTimes{
		RestrictionType:     "Daily",
		DayRestrictionTimes: &StartAndEndTime{StartTime: "20:00", EndTime: "04:00"},
	})
	value := data["restrictionTimes"].(map[string]interface{})["value"].(map[string]interface{})
	day := value["dayRestrictionTimes"].(map[string]interface{})
	if day["startTime"] != "20:00" || day["endTime"] != "04:00" {
		t.Errorf("day restriction wrong: %v", day)
	}
}

func TestAttachRestrictionTimes_Nil(t *testing.T) {
	t.Parallel()
	data := map[string]interface{}{}
	attachRestrictionTimes(data, "restrictionTimes", nil)
	if _, present := data["restrictionTimes"]; present {
		t.Error("nil input should leave key absent")
	}
}
