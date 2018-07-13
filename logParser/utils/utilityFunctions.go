package utils

import (
	"sort"
	"time"
)

// DataRow represents a single row containing a request time and URL
type DataRow struct {
	RequestTime time.Time
	RequestURL  string
}

// RequestStatistics represents a row containing a request URL,
// total hist and distribution
type RequestStatistics struct {
	RequestURL   string
	Times        int
	Distribution float32
	MaxTPM       int
}

// CountTimesPerRequestAndDistribution calculates Hits and Distribution for each request
func CountTimesPerRequestAndDistribution(data []DataRow, sliceTo int) ([]RequestStatistics, int) {
	totalHits := float32(len(data))
	tmpMapRequestStats := make(map[string]int)
	tmpMapHitsPerMinStats := make(map[time.Time]map[string]int)

	for _, row := range data {
		if _, ok := tmpMapRequestStats[row.RequestURL]; ok {
			tmpMapRequestStats[row.RequestURL]++
		} else {
			tmpMapRequestStats[row.RequestURL] = 1
		}

		if perMin, ok := tmpMapHitsPerMinStats[row.RequestTime]; ok {
			if _, ok := perMin[row.RequestURL]; ok {
				perMin[row.RequestURL]++
			} else {
				perMin[row.RequestURL] = 1
			}
		} else {
			tmpMapHitsPerMinStats[row.RequestTime] = map[string]int{row.RequestURL: 1}
		}
	}

	maxTPMPerURL, maxTPM := findMaxTPMPerURL(tmpMapHitsPerMinStats)

	times := make([]RequestStatistics, len(tmpMapRequestStats))
	i := 0
	for key, val := range tmpMapRequestStats {
		times[i] = RequestStatistics{
			RequestURL:   key,
			Times:        val,
			Distribution: (float32(val) / totalHits) * 100.0,
			MaxTPM:       maxTPMPerURL[key],
		}
		i++
	}

	sort.Slice(times, func(i, j int) bool { return times[i].Times > times[j].Times })

	return times[:sliceTo], maxTPM
}

func findMaxTPMPerURL(data map[time.Time]map[string]int) (map[string]int, int) {
	parsed := make(map[string]int)
	maxTPM := 0
	for _, val := range data {
		sum := 0
		for url, times := range val {
			sum += times
			if curMax, ok := parsed[url]; ok {
				if times > curMax {
					parsed[url] = times
				}
			} else {
				parsed[url] = times
			}
		}
		if sum > maxTPM {
			maxTPM = sum
		}
	}

	return parsed, maxTPM
}
