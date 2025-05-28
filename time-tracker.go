package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

type Activity struct {
	Name     string
	Duration time.Duration
	Date     time.Time
}

const fileName = "activity_log.csv"

func main() {
	fmt.Println("Welcome to the Time Tracker CLI!")
	fmt.Println("Commands:")
	fmt.Println("  start [activity_name] - Start tracking an activity.")
	fmt.Println("  quit                  - Exit the application.")

	for {
		fmt.Print("\n> ")
		var input string
		fmt.Scanln(&input)

		if strings.ToLower(input) == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		if strings.ToLower(input) == "start" {
			fmt.Print("Enter the activity name: ")
			var activityName string
			fmt.Scanln(&activityName)

			startTime := time.Now()
			fmt.Printf("Started tracking: %s. Type 'stop' to stop tracking.\n", activityName)

			stopChannel := make(chan bool)
			go showLiveTime(startTime, stopChannel)

			for {
				var stopCommand string
				fmt.Scanln(&stopCommand)

				if strings.ToLower(stopCommand) == "stop" {
					stopChannel <- true
					close(stopChannel)

					duration := time.Since(startTime)
					activity := Activity{
						Name:     activityName,
						Duration: duration,
						Date:     startTime,
					}

					saveActivity(activity)
					fmt.Printf("Tracked activity '%s' for %s on %s.\n", activity.Name, activity.Duration, activity.Date.Format("2006-01-02"))
					break
				}
			}
		} else {
			fmt.Println("Unknown command. Please use 'start [activity_name]' or 'quit'.")
		}
	}
}

func showLiveTime(startTime time.Time, stopChannel chan bool) {
	for {
		select {
		case <-stopChannel:
			return
		default:
			elapsed := time.Since(startTime)
			fmt.Printf("\rTracking... Elapsed Time: %s", elapsed.Truncate(time.Second))
			time.Sleep(1 * time.Second)
		}
	}
}

func saveActivity(activity Activity) {
	// Check if the file exists
	_, err := os.Stat(fileName)
	isNewFile := os.IsNotExist(err)

	// Open the file for appending, creating it if necessary
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers if the file is new
	if isNewFile {
		if err := writer.Write([]string{"Activity Name", "Duration", "Date"}); err != nil {
			fmt.Printf("Error writing headers: %v\n", err)
			return
		}
	}

	// Write the activity details
	if err := writer.Write([]string{
		activity.Name,
		activity.Duration.String(),
		activity.Date.Format("2006-01-02 15:04:05"),
	}); err != nil {
		fmt.Printf("Error writing activity: %v\n", err)
	}
}
