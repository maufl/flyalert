// Copyright (C) 2014 Felix Maurer
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>

package main

import (
	"fmt"
	"time"
	forecast "github.com/mlbright/forecast/v2"
)

type Range struct {
	Max float64
	Min float64
}

func (r Range) String() string {
	return "[" + fmt.Sprint(r.Min) + "-" + fmt.Sprint(r.Max) + "]"
}

type Trigger struct {
	Long string
	Lat string
	CloudCover Range
	PrecipProbability Range
	PrecipIntensity Range
	WindBearing Range
	WindSpeed Range
}

func (t Trigger) Notification(time time.Time) Notification {
	return Notification{Day: time, Long: t.Long, Lat: t.Lat}
}

func (t Trigger) Process() (notifications Notifications) {
	f, err := forecast.Get(conf.ApiKey, t.Lat, t.Long, "now", forecast.CA)
	if err != nil {
		fmt.Println("Could not get forcast")
		return
	}
	fmt.Println("Timezone ", f.Timezone, " offset ", f.Offset)
	for _, dataPoint := range f.Daily.Data {
		fmt.Println()
		time := time.Unix(int64(dataPoint.Time), 0)
		fmt.Print("T ", time, " | CC ", dataPoint.CloudCover, " | PP ", dataPoint.PrecipProbability,)
		fmt.Println(" | PI ", dataPoint.PrecipIntensity, " | WB ", dataPoint.WindBearing, " | WS ", dataPoint.WindSpeed)

		if t.CloudCover.Min > dataPoint.CloudCover || t.CloudCover.Max < dataPoint.CloudCover {
			fmt.Println("CouldCover is not in range, trigger failed")
			continue
		}
		if t.PrecipIntensity.Min > dataPoint.PrecipIntensity || t.PrecipIntensity.Max < dataPoint.PrecipIntensity {
			fmt.Println("Precipitation intensity is not in range, trigger failed")
			continue
		}
		if t.PrecipProbability.Min > dataPoint.PrecipProbability || t.PrecipProbability.Max < dataPoint.PrecipProbability {
			fmt.Println("Precipitation probability is not in range, trigger failed")
			continue
		}
		if t.WindBearing.Min > dataPoint.WindBearing || t.WindBearing.Max < dataPoint.WindBearing {
			fmt.Println("Wind bearing is not in range, trigger failed")
			continue
		}
		if t.WindSpeed.Min > dataPoint.WindSpeed || t.WindSpeed.Max < dataPoint.WindSpeed {
			fmt.Println("Wind speed is not in range, trigger failed")
			continue
		}
		fmt.Println("Trigger succeeded for day ", time)
		notifications = append(notifications, t.Notification(time))
	}
	return
}
