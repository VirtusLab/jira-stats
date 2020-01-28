package domain

import (
	"log"
	"time"
)

const StatusField = "status"
const StateDev = "In Development"

type Now func() time.Time

type DaysCalculator struct {
	ClockNow Now
}

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
func (this *DaysCalculator) CalculateDevDays(ticket Ticket, start time.Time, end time.Time) float64 {
	var cumulativeTime int

	if this.shouldSkipTicket(ticket) {
		log.Printf("Ticket %s has been skipped from dev time calculation", ticket.Key)
		return 0.0
	}

	intervals := MakeIntervals(ticket)

	for _, interval := range intervals {
		devTime := this.calculateDevTime(interval, start, end)
		cumulativeTime += devTime
	}

	return float64(cumulativeTime) / 8.0
}

func (this *DaysCalculator) shouldSkipTicket(ticket Ticket) bool {
	return ticket.Type == "Epic" // epics are being skipped from calculation
}

func (this *DaysCalculator) calculateDevTime(interval TransitionInterval, start time.Time, end time.Time) int {
	//log.Printf(interval.ToString())

	// we are not interested in non-dev states or state intervals outside given boundaries
	if interval.State != StateDev || !this.isTransitionRelevantForBoundaries(interval, start, end) {
		return 0
	}

	interval = this.adjustDatesToBounds(interval, start, end)

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
		return this.calculateWorkingHours(interval.Start, interval.End)
		return int(diff.Hours())
	}
}

func (this *DaysCalculator) calculateWorkingHours(start time.Time, end time.Time) int {
	totalHours := 0

	if this.isWorkingDay(start) { // take care of first day
		if start.Hour() < 12 {
			totalHours += 8
		} else {
			totalHours += 4
		}
	}

	if this.isWorkingDay(end) { // take care of last day
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

		if this.isWorkingDay(currentDay) {
			totalHours += 8
		}

		currentDay = currentDay.Add(time.Hour * 24)
	}

	return totalHours
}

func (this *DaysCalculator) adjustDatesToBounds(interval TransitionInterval, start time.Time, end time.Time) TransitionInterval {
	if interval.Start.Before(start) {
		interval.Start = start
	}

	endBound := this.calculateEndBound(end)
	if endBound.Before(interval.End) {
		interval.End = endBound
	}

	return interval
}

func (this *DaysCalculator) calculateEndBound(end time.Time) time.Time {
	now := this.now()
	endInterval := end
	if now.Before(end) {
		endInterval = now
	}
	return endInterval
}

func (this *DaysCalculator) isWorkingDay(date time.Time) bool {
	return date.Weekday() != time.Saturday && date.Weekday() != time.Sunday
}

func (this *DaysCalculator) isTransitionRelevantForBoundaries(interval TransitionInterval, start time.Time, end time.Time) bool {
	return (interval.Start.Before(start) && interval.End.After(end)) || // interval contains boundaries
		(this.between(interval.Start, start, end) || this.between(interval.End, start, end)) // intervals overlaps or is contained withing boundaries

}

func (this *DaysCalculator) between(date time.Time, lowBound time.Time, highBound time.Time) bool {
	return date.After(lowBound) && date.Before(highBound)
}

func (this *DaysCalculator) now() time.Time {
	if this.ClockNow != nil {
		return this.ClockNow()
	} else {
		return time.Now()
	}
}
