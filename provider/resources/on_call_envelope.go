package resources

// OneUptime's on-call schedule layer stores two structured JSON columns —
// rotation (Recurring) and restrictionTimes (RestrictionTimes) — each in the
// same {_type, value} envelope that monitorSteps uses. The helpers here
// mirror the monitor_steps_envelope.go pattern scoped to those two types.

// attachRotation writes a Recurring envelope to data[key]. A nil pointer is
// a no-op so callers can conditionally set only the fields they want to
// change on Update.
func attachRotation(data map[string]interface{}, key string, r *Recurring) {
	if r == nil {
		return
	}
	data[key] = map[string]interface{}{
		"_type": "Recurring",
		"value": map[string]interface{}{
			"intervalType": r.IntervalType,
			// intervalCount is itself a PositiveNumber envelope on the wire —
			// OneUptime's fromJSON accepts raw numbers too, but emitting the
			// envelope matches what MonitorSteps-style consumers produce and
			// avoids any server-side DatabaseProperty parse quirks.
			"intervalCount": map[string]interface{}{
				"_type": "PositiveNumber",
				"value": r.IntervalCount,
			},
		},
	}
}

// attachRestrictionTimes writes a RestrictionTimes envelope to data[key].
// The upstream TypeScript uses the misspelled key `restictionType` (missing
// `r`) and we preserve that on the wire even though the Pulumi input uses
// the correct spelling `restrictionType`.
func attachRestrictionTimes(data map[string]interface{}, key string, rt *RestrictionTimes) {
	if rt == nil {
		return
	}
	value := map[string]interface{}{
		"restictionType":         rt.RestrictionType, // preserve upstream typo
		"weeklyRestrictionTimes": weeklyRestrictionsToWire(rt.WeeklyRestrictionTimes),
	}
	if rt.DayRestrictionTimes != nil {
		value["dayRestrictionTimes"] = map[string]interface{}{
			"startTime": rt.DayRestrictionTimes.StartTime,
			"endTime":   rt.DayRestrictionTimes.EndTime,
		}
	} else {
		value["dayRestrictionTimes"] = nil
	}
	data[key] = map[string]interface{}{
		"_type": "RestrictionTimes",
		"value": value,
	}
}

func weeklyRestrictionsToWire(ws []WeeklyRestriction) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ws))
	for _, w := range ws {
		out = append(out, map[string]interface{}{
			"startDay":  w.StartDay,
			"endDay":    w.EndDay,
			"startTime": w.StartTime,
			"endTime":   w.EndTime,
		})
	}
	return out
}
