package domain

import (
	"time"
)

const StateDev = "In Development"

/**
	Calculating length of development is non-trivial. Mostly because we don't know how much really someone worked
	during given day, we only know discrete state transition points.

	Following rules therefore has been assumed:
	- day is 8h of work
	- if DEV started ticket before noon - whole day is going to be counted for that day
	- if DEV started ticket after noon - day would be counted as half (0.5)
    - if DEV started and moved ticket further at the same day - number of ours would be rounded to nearest multiplication of 0.25 of day (2 hours)
	- weekends are not counted
*/
func CalculateDevDays(ticket Ticket, start time.Time, end time.Time) float64 {
	var cumulativeTime int

	for _, transition := range ticket.Transitions {
		devTime := calculateDevTime(transition, start, end)
		cumulativeTime += devTime
	}

	return float64(cumulativeTime) / 8.0
}

func calculateDevTime(interval TransitionInterval, start time.Time, end time.Time) int {
	//log.Printf(interval.ToString())

	// we are not interested in non-dev states or state intervals outside given boundaries
	if interval.State != StateDev || !isTransitionRelevantForBoundaries(interval, start, end) {
		return 0
	}

	interval = adjustDatesToBounds(interval, start, end)

	diff := interval.End.Sub(interval.Start)

	if diff.Minutes() <= 2*60 {
		return 2
	} else if diff.Minutes() <= 4*60 {
		return 4
	} else if diff.Minutes() <= 6*60 {
		return 6
	} else if diff.Minutes() <= 8*60 {
		return 8
	} else {
		return calculateWorkingHours(interval.Start, interval.End)
		return int(diff.Hours())
	}
}

func calculateWorkingHours(start time.Time, end time.Time) int {
	totalHours := 0

	if isWorkingDay(start) { // take care of first day
		if start.Hour() < 12 {
			totalHours += 8
		} else {
			totalHours += 4
		}
	}

	if isWorkingDay(end) { // take care of last day
		if end.Hour() < 12 {
			totalHours += 4
		} else {
			totalHours += 8
		}
	}

	// calculate all days in between
	currentDay := start
	currentDay = currentDay.Add(time.Hour * 24)

	for currentDay.Year() <= end.Year() && currentDay.YearDay() < end.YearDay() {

		if isWorkingDay(currentDay) {
			totalHours += 8
		}

		currentDay = currentDay.Add(time.Hour * 24)
	}

	//if (t.Weekday() != 6 && t.Weekday() != 7) {
	//	days++
	//}
	return totalHours
}

func adjustDatesToBounds(interval TransitionInterval, start time.Time, end time.Time) TransitionInterval {
	if interval.Start.Before(start) {
		interval.Start = start
	}

	if interval.End.After(end) {
		interval.End = end
	}

	return interval
}

func isWorkingDay(date time.Time) bool {
	return date.Weekday() != time.Saturday && date.Weekday() != time.Sunday
}

func isTransitionRelevantForBoundaries(interval TransitionInterval, start time.Time, end time.Time) bool {
	return (interval.Start.Before(start) && interval.End.After(end)) || // interval contains boundaries
		(between(interval.Start, start, end) || between(interval.End, start, end)) // intervals overlaps or is contained withing boundaries

}

func between(date time.Time, lowBound time.Time, highBound time.Time) bool {
	return date.After(lowBound) && date.Before(highBound)
}
