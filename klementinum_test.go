package klementinum

import (
	"github.com/chlachula/klementinum/data"
	"testing"
)

func TestTemperatureStatistics(t *testing.T) {
	var testData = []data.TempRecord{
		{Y: 1775, M: 1, D: 1, T: -7.0},
		{Y: 1775, M: 1, D: 2, T: -2.2},
		{Y: 1775, M: 1, D: 3, T: -1.0},
		{Y: 1775, M: 1, D: 4, T: 0.1},
		{Y: 1775, M: 1, D: 5, T: 2.2},
		{Y: 1775, M: 1, D: 6, T: 3.2},
		{Y: 1775, M: 1, D: 7, T: 3.5},
		{Y: 1775, M: 1, D: 8, T: 4.1},
		{Y: 1775, M: 1, D: 9, T: 4.0},
	}
	var want TStat
	want.MinT = data.TempRecord{Y: 1775, M: 1, D: 1, T: -7.0}
	want.MaxT = data.TempRecord{Y: 1775, M: 1, D: 8, T: 4.1}

	if got := TemperatureStatistics(testData); got != want {
		t.Errorf("Message() = %v, want %v", got, want)
	}

}
func DaysInYearTest(t *testing.T) {
	a := DaysInYear(1775) == 365
	b := DaysInYear(1968) == 366
	c := DaysInYear(2000) == 365
	d := DaysInYear(2020) == 366
	if !(a && b && c && d){
		t.Errorf("DaysInYear() returned wrong days 365 or 366")
	}
}

