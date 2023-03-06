/* Example how to create, read and remove a PKLM sample csv file
 */
package klementinum

import (
	"fmt"

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

func svgAverage() string {
	svgFormat := `<svg
 xmlns="http://www.w3.org/2000/svg" 
 xmlns:xlink="http://www.w3.org/1999/xlink" 
 viewBox="0 0 600 300" 
 width="600" height="300">
 <title>Yearly differences of the Klementinum temperature records</title>
 <defs>
	<style>
	   svg { background-color: lightgray; }
	   text { font-size: 0.9px; }
	</style>
	<pattern id="bg_image" patternUnits="userSpaceOnUse" width="500" height="500">
	   <image href="SolRootRiverPark1.jpg" x="0" y="0" width="500" height="500" />
	</pattern>
 </defs>
 <g id="main">
    <rect x="0" y="0" width="100" height="300" fill="url(#bg_image)" stroke="none" />
    <text x="10" y="10" >Average year temperature: %.1f°C</text>
	<path d="M0,150 l0,600" stroke="black" stroke-width="10" />
 </g>
 %s
</svg>
`
	return fmt.Sprintf(svgFormat, 9.9)
}
