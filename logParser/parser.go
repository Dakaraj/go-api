package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	re "regexp"
	"strings"
	"time"

	"github.com/dakaraj/go-api/logParser/utils"
)

const (
	// Use limitTo variable to limit results only to N top requests
	limitTo = 25
	// Use filePath variable to set a custom path to the log file
	filePath = "logParser/data/access.log"
	// Use includeStatic variable to include static resources to statistics
	includeStatic = true

	// List of pattern strings to use later for URL parsing and truncation
	logLinePatternString            = `\[(\d{2}/\w{3}/\d{4}:\d{2}:\d{2}):\d{2} .*? "((?:GET|POST|HEAD|OPTIONS) \S*)`
	parseFilePatternTruncateString  = `/servlets/\d+?Dispatch/\d+?/jspforward\?file`
	fileToPageTruncatePatternString = `\?file=.+?&page=`
	pageTruncatePatternString       = `&page=([^&#=]+).*`
	staticResoursesPatternString    = `\.(:?jpg|gif|png|js|css|swf|woff)$`
)

var totalRequests int

func main() {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer file.Close()

	logLinePattern := re.MustCompile(logLinePatternString)
	staticResoursesPattern := re.MustCompile(staticResoursesPatternString)
	parseFilePatternTruncate := re.MustCompile(parseFilePatternTruncateString)
	fileToPageTruncatePattern := re.MustCompile(fileToPageTruncatePatternString)
	pageTruncatePattern := re.MustCompile(pageTruncatePatternString)

	reader := bufio.NewReader(file)

	var parsed []utils.DataRow

	for {
		bytesLine, _, err := reader.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
		}
		line := string(bytesLine)
		results := logLinePattern.FindStringSubmatch(line)
		if len(results) != 3 {
			continue
		}
		timeString, rawURL := results[1], results[2]
		timeVal, _ := time.Parse("2/Jan/2006:15:04", timeString)

		if staticResoursesPattern.MatchString(rawURL) {
			if includeStatic {
				rawURL = "<<staticResourceRequest>>"
			} else {
				continue
			}
		} else if matched, _ := re.MatchString(`\?file=`, rawURL); matched {
			rawURL = parseFilePatternTruncate.ReplaceAllString(rawURL, "/servlets/***Dispatch/***/jspforward?file")
			if strings.Contains(rawURL, "&page=") {
				rawURL = fileToPageTruncatePattern.ReplaceAllString(rawURL, "?file=***&page=")
				rawURL = pageTruncatePattern.ReplaceAllString(rawURL, "&page=${1}")
			}
		}

		parsed = append(parsed, utils.DataRow{RequestTime: timeVal, RequestURL: rawURL})
	}

	totalRequests = len(parsed)
	fmt.Printf("Total reqests: %d\n", totalRequests)
	timesPerRequestAndDistribution, maxTPM := utils.CountTimesPerRequestAndDistribution(parsed, limitTo)
	fmt.Printf("Max TPM: %d\n", maxTPM)
	for _, tpr := range timesPerRequestAndDistribution {
		fmt.Printf("Request: %s, Times: %d, Distribution: %.2f%%, Max TPM: %d\n",
			tpr.RequestURL, tpr.Times, tpr.Distribution, tpr.MaxTPM)
	}
}
