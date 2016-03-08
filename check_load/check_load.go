package main

import (
	"fmt"
	"strings"
	"strconv"

	"github.com/olorin/nagiosplugin"
	"github.com/shirou/gopsutil/load"
	"gopkg.in/alecthomas/kingpin.v2"
)


type WarningLimit struct {
	Load1  float64
	Load5  float64
	Load15 float64
}
type CriticalLimit struct {
	Load1  float64
	Load5  float64
	Load15 float64
}
type Limits struct {
	Warning  WarningLimit
	Critical CriticalLimit
}


var (
  debug   = kingpin.Flag("debug", "Enable debug mode.").Bool()

	warning  = kingpin.Flag("warning", "The warning value to trigger an alert").Required().Short('w').String()
	critical  = kingpin.Flag("critical", "The critical value to trigger an alert").Required().Short('c').String()

)



func parseThreshholds() (l Limits) {
	w := strings.Split(*warning, ",")
	c := strings.Split(*critical, ",")


	for {
		if len(w) >= 3 {
			break
		}
		w = append(w, w[len(w)-1])
	}
	for {
		if len(c) >= 3 {
			break
		}
		c = append(c, c[len(c)-1])
	}

	w1, _ := strconv.ParseFloat(w[0], 64)
	w5, _ := strconv.ParseFloat(w[1], 64)
	w15, _ := strconv.ParseFloat(w[2], 64)
	c1, _ := strconv.ParseFloat(c[0], 64)
	c5, _ := strconv.ParseFloat(c[1], 64)
	c15, _ := strconv.ParseFloat(c[2], 64)

	l = Limits {
    Warning: WarningLimit {
        Load1:  w1,
        Load5:  w5,
				Load15: w15,
    },
		Critical: CriticalLimit {
        Load1:  c1,
        Load5:  c5,
				Load15: c15,
    },
	}

	return
}

func CheckLoadAvg() {
  kingpin.Parse()
	kingpin.Version("2016.03.01")

	check := nagiosplugin.NewCheck()
	defer check.Finish()

	limits := parseThreshholds()

	v, err := load.LoadAvg()
	if err != nil {
		check.AddResult(nagiosplugin.WARNING, "Error getting load average")
	}

	check.AddPerfDatum("loadavg1", "", v.Load1, 0, 100, limits.Warning.Load1, limits.Critical.Load1)
	check.AddPerfDatum("loadavg5", "", v.Load5, 0, 100, limits.Warning.Load5, limits.Critical.Load5)
	check.AddPerfDatum("loadavg15", "", v.Load15, 0, 100, limits.Warning.Load15, limits.Critical.Load15)

	if v.Load1 > limits.Critical.Load1 {
		check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("Load Average min 1 is currently %v", v.Load1))
	} else if v.Load1 > limits.Warning.Load1 {
		check.AddResult(nagiosplugin.WARNING, fmt.Sprintf("Load Average min 1 is currently %v", v.Load1))
	}

	if v.Load5 > limits.Critical.Load5 {
		check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("Load Average min 5 is currently %v", v.Load5))
	} else if v.Load5 > limits.Warning.Load5 {
		check.AddResult(nagiosplugin.WARNING, fmt.Sprintf("Load Average min 5 is currently %v", v.Load5))
	}

	if v.Load15 > limits.Critical.Load15 {
		check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("Load Average min 15 is currently %v", v.Load15))
	} else if v.Load15 > limits.Warning.Load15 {
		check.AddResult(nagiosplugin.WARNING, fmt.Sprintf("Load Average min 15 is currently %v", v.Load15))
	}

	// everything is good
	check.AddResult(nagiosplugin.OK, "Load Average is under the limits")
}


func main() {
	CheckLoadAvg()
}
