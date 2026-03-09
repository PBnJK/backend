package converter

/* Conversion rate from gram to other units
 * Sources:
 *   https://en.wikipedia.org/wiki/Gram#Conversion_factors
 *   WolframAlpha
 */
var toGramMap = map[string]float64{
	"w-g":   1,
	"w-kg":  1000,
	"w-oz":  28.3,
	"w-lbs": 453.59237,
}

// Converts a number expressed in the weight unit "from" to unit "to"
func ConvertWeight(number float64, from string, to string) float64 {
	if from == to {
		return number
	}

	inGrams := number * toGramMap[from]
	return inGrams / toGramMap[to]
}
