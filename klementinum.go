/* Example how to create, read and remove a PKLM sample csv file
 */
package klementinum

import (
	"github.com/chlachula/klementinum/data"
)

//all := data.TemperatureRecords()

type TStat = struct {
	MinT data.TempRecord
	MaxT data.TempRecord
}

func TemperatureStatistics(allData []data.TempRecord) TStat {
	var s TStat
	s.MinT.T = 100.0
	s.MaxT.T = -273.0

	for _, r := range allData {
		if r.T < s.MinT.T {
			s.MinT = r
		}
		if r.T > s.MaxT.T {
			s.MaxT = r
		}
	}
	return s
}
