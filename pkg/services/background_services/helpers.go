package background_services

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

var lastTickTime = time.Now()

func isScientificNotation(str string) bool {
	r := regexp.MustCompile(`^[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?$`)
	return r.MatchString(str)
}

// compareFiles compares two file names based on month and year
// compareFiles compares two file names based on month and year
func compareFiles(file1, file2 string) bool {
	date1, err := dateFromFile(file1)
	if err != nil {
		log.Println("Error parsing date for", file1, ":", err)
		return false
	}
	date2, err := dateFromFile(file2)
	if err != nil {
		log.Println("Error parsing date for", file2, ":", err)
		return false
	}

	return date1.Before(date2)
}

// getMonthYearWithDay extracts the month and year from the file name and returns it as a time.Time object
// dateFromFile extracts the date from the file name and returns it as a time.Time object
func dateFromFile(file string) (time.Time, error) {
	// Define a regular expression pattern to match the date part in the file name
	datePattern := `[a-zA-Z]{3}\d{1,2}_\d{4}`
	regex := regexp.MustCompile(datePattern)

	// Find the date part in the file name using the regular expression
	dateString := regex.FindString(file)
	if dateString == "" {
		return time.Time{}, fmt.Errorf("invalid file name format: no date found")
	}

	// Parse the date using the known format
	t, err := time.Parse("Jan2_2006", dateString)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func convertScientificNotation(appID string) (int, error) {
	parsedAppID, err := strconv.ParseFloat(appID, 64)
	if err != nil {
		return 0, err
	}
	appIDInt := int(parsedAppID)

	return appIDInt, nil
}

var updateInterval time.Duration // Define the update interval as a package-level variable

// SetUpdateInterval sets the update interval used by the background process
func SetUpdateInterval(interval time.Duration) {
	updateInterval = interval
}

// GetUpdateInterval returns the update interval used by the background process
func GetUpdateInterval() time.Duration {
	return updateInterval
}

func GetLastTickTime() time.Time {
	return lastTickTime
}

func SetLastTickTime(t time.Time) {
	lastTickTime = t
}

func TimeRemainingForNextUpdate(updateInterval time.Duration) time.Duration {

	lastProcessTime := GetLastTickTime()
	now := time.Now()
	nextUpdate := lastProcessTime.Add(updateInterval)
	if now.Before(nextUpdate) {
		return nextUpdate.Sub(now)
	}

	return -1
}

func TimeInDaysFormat(duration time.Duration) string {
	if duration < 0 {
		return "N/A"
	}

	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	return fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds", days, hours, minutes, seconds)
}
