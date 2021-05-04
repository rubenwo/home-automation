package models

type DayOfWeek string

const (
	Star      DayOfWeek = "*"
	Monday              = "Monday"
	Tuesday             = "Tuesday"
	Wednesday           = "Wednesday"
	Thursday            = "Thursday"
	Friday              = "Friday"
	Saturday            = "Saturday"
	Sunday              = "Sunday"
)

func (d DayOfWeek) Int() int {
	switch d {
	case Monday:
		return 1
	case Tuesday:
		return 2
	case Wednesday:
		return 3
	case Thursday:
		return 4
	case Friday:
		return 5
	case Saturday:
		return 6
	case Sunday:
		return 7
	}

	return -1
}
