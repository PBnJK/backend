package converter

/* Conversion rate from meter to other units
 * Sources:
 *   https://en.wikipedia.org/wiki/Metre#Equivalents_in_other_units
 *   WolframAlpha
 */
var toMeterMap = map[string]float64{
	"l-mm":  0.001,
	"l-cm":  0.01,
	"l-dm":  0.1,
	"l-m":   1,
	"l-dam": 10,
	"l-hm":  100,
	"l-km":  1000,
	"l-in":  0.0254,
	"l-ft":  0.3048,
	"l-yd":  0.9144,
	"l-mi":  1609.344,
}

// Converts a number expressed in the length unit "from" to unit "to"
func ConvertLength(number float64, from string, to string) float64 {
	if from == to {
		return number
	}

	inMeters := number * toMeterMap[from]
	return inMeters / toMeterMap[to]
}
