/* Example how to create, read and remove a PKLM sample csv file
 */
package klementinum

import (
	"github.com/chlachula/klementinum/data"
)

//all := data.TemperatureRecords()

type TStat = struct {
	Year1    int
	YearEnd  int
	MinT     data.TempRecord
	MaxT     data.TempRecord
	YearTavg []float32
}

func DaysInYear(y int) int {
	yDays := 365
	if y%4 == 0 && y != 2000 {
		yDays = 366
	}
	return yDays
}
func TemperatureStatistics(allData []data.TempRecord) TStat {
	var s TStat
	s.Year1 = allData[0].Y
	s.YearEnd = allData[len(allData)-1].Y
	s.YearTavg = make([]float32, s.YearEnd-s.Year1+1)
	s.MinT.T = 100.0
	s.MaxT.T = -273.0

	for _, r := range allData {
		s.YearTavg[r.Y-s.Year1] += r.T
		if r.T < s.MinT.T {
			s.MinT = r
		}
		if r.T > s.MaxT.T {
			s.MaxT = r
		}
	}
	for y, _ := range s.YearTavg {
		year := y + s.Year1
		s.YearTavg[y] = s.YearTavg[y] / float32(DaysInYear(year))
	}

	return s
}
