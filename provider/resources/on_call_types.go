package resources

// Shared types for on-call resources. Like monitor_types.go, these mirror
// OneUptime's TypeScript interfaces (master as of 2026-04-20) and are
// wrapped into `{_type, value}` envelopes by the helpers below on write.

// Recurring describes a rotation cadence used by on-call schedule layers
// (and eventually by scheduled maintenance templates). Serializes on the
// wire as `{_type: "Recurring", value: {intervalType, intervalCount}}`.
type Recurring struct {
	// IntervalType is one of "Hour" | "Day" | "Week" | "Month" | "Year".
	IntervalType  string `pulumi:"intervalType" json:"intervalType"`
	IntervalCount int    `pulumi:"intervalCount" json:"intervalCount"`
}

// RestrictionTimes limits when a schedule layer is active. Serializes as
// `{_type: "RestrictionTimes", value: {restictionType [sic], ...}}`.
//
// Note: the upstream TypeScript has a typo — `restictionType` (no second r)
// — that is preserved on the wire. The Pulumi surface uses the correct
// spelling (`restrictionType`) and the envelope helper rewrites the key.
type RestrictionTimes struct {
	// RestrictionType is "None" | "Daily" | "Weekly".
	RestrictionType        string                `pulumi:"restrictionType" json:"-"`
	DayRestrictionTimes    *StartAndEndTime      `pulumi:"dayRestrictionTimes,optional" json:"-"`
	WeeklyRestrictionTimes []WeeklyRestriction   `pulumi:"weeklyRestrictionTimes,optional" json:"-"`
}

type StartAndEndTime struct {
	StartTime string `pulumi:"startTime" json:"startTime"`
	EndTime   string `pulumi:"endTime" json:"endTime"`
}

type WeeklyRestriction struct {
	// StartDay / EndDay use OneUptime's day names: "Sunday", "Monday", ....
	StartDay  string `pulumi:"startDay" json:"startDay"`
	EndDay    string `pulumi:"endDay" json:"endDay"`
	StartTime string `pulumi:"startTime" json:"startTime"`
	EndTime   string `pulumi:"endTime" json:"endTime"`
}
