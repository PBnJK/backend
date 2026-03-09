package converter

import "log"

// Converts a number expressed in the temperature unit "from" to unit "to"
func ConvertTemperature(number float64, from string, to string) float64 {
	if from == to {
		return number
	}

	switch from {
	case "t-C":
		if to == "t-F" {
			return number*9/5.0 + 32
		} else {
			return number + 273
		}
	case "t-F":
		if to == "t-C" {
			return (number - 32) / (9 / 5.0)
		} else {
			return (number + 459.67) * 5 / 9.0
		}
	case "t-K":
		if to == "t-C" {
			return number - 273
		} else {
			return number*9/5.0 - 459.67
		}
	}

	log.Fatal("unreachable")
	return 0
}
