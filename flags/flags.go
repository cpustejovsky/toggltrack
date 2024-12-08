package flags

import (
	"flag"
	"strconv"
)

func Float64(name, envVar string, def float64, description string) *float64 {
	return flag.Float64(name, envVarToFloat64(envVar, def), description)
}

func Int(name, envVar string, def int, description string) *int {
	return flag.Int(name, envVarToInt(envVar, def), description)
}

func Bool(name, envVar string, def bool, description string) *bool {
	return flag.Bool(name, envVarToBool(envVar, def), description)
}
func envVarToFloat64(s string, def float64) float64 {
	v, err := strconv.ParseFloat(s, 0)
	if err != nil {
		return def
	}
	return v
}

func envVarToInt(s string, def int) int {
	v, err := strconv.ParseInt(s, 0, 0)
	if err != nil {
		return def
	}
	return int(v)
}

func envVarToBool(s string, def bool) bool {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return def
	}
	return v
}
