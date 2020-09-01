package main

func MaxUserJourneyResponseLatency(r []UserJourneyResult) int64{
	// GetMaxUserJourneyLatency is a custom max function
	// used to get the biggest latency value in a set of UserJourneys.

	// I think it's nicer to get actual data from r rather
	// than just initialise this to some random number.
	var max int64 = r[0].JourneyResponseTimeMS

	for _, v := range r[1:]{
		if v.JourneyResponseTimeMS > max{
			max = v.JourneyResponseTimeMS
		}
	}

	return max

}

func MinUserJourneyResponseLatency(r []UserJourneyResult) int64{
	min := r[0].JourneyResponseTimeMS

	for _, v := range r[1:]{
		if v.JourneyResponseTimeMS < min{
			min = v.JourneyResponseTimeMS
		}
	}
	return min
}