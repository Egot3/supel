package types

import ttpb "github.com/Egot3/supel/backend/contracts/timetable"

type Day string

const (
	MONDAY    = "Monday"
	TUESDAY   = "Tuesday"
	WEDNESDAY = "Wednesday"
	THURSDAY  = "Thursday"
	FRIDAY    = "Friday"
	SATURDAY  = "Saturday"
	SUNDAY    = "Sunday"
)

func ProtoDayToDay(day ttpb.Day) Day {
	switch day {
	case ttpb.Day_MONDAY:
		return MONDAY
	case ttpb.Day_TUESDAY:
		return TUESDAY
	case ttpb.Day_WEDNESDAY:
		return WEDNESDAY
	case ttpb.Day_THURSDAY:
		return THURSDAY
	case ttpb.Day_FRIDAY:
		return FRIDAY
	case ttpb.Day_SATURDAY:
		return SATURDAY
	case ttpb.Day_SUNDAY:
		return SUNDAY
	default:
		return MONDAY
	}
}
