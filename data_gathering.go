package main

func MaxUserJourneyResponseLatency(r []UserJourneyResult) int64 {
	// GetMaxUserJourneyLatency is a custom max function
	// used to get the biggest latency value in a set of UserJourneys.

	// I think it's nicer to get actual data from r rather
	// than just initialise this to some random number.
	var max int64 = r[0].JourneyResponseTimeMS

	for _, v := range r[1:] {
		if v.JourneyResponseTimeMS > max {
			max = v.JourneyResponseTimeMS
		}
	}
	return max
}

func MinUserJourneyResponseLatency(r []UserJourneyResult) int64 {
	min := r[0].JourneyResponseTimeMS

	for _, v := range r[1:] {
		if v.JourneyResponseTimeMS < min {
			min = v.JourneyResponseTimeMS
		}
	}
	return min
}

//func FindPercentile(latencies []int, p int) int{
//	// FindPercentile looks for the according percentile
//	// figure depending on the given p.
//
//	// Check everything is sorted.
//	if !sort.IntsAreSorted(latencies){
//		sort.Ints(latencies)
//	}
//
//	// 1: percentile/100 * number of items in list
//	fmt.Println(float64(p)/100)
//
//
//	return 100
//
//}

func CountResponseCodes(r []UserJourneyResult) map[int]int {
	// CountResponseCodes is here to go through a load of responses
	// and to generate a single, flat map that shows, of all the
	// URLs visited, what response codes we received. It can be used
	// later to help with data analysis, or just to get a quick view of
	// what responses came up most often in a test.

	codes := map[int]int{}
	for _, v := range r {
		for code, count := range v.Codes {
			codes[code] += count
		}
	}
	return codes
}
