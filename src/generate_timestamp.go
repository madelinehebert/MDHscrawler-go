package main

import (
	"fmt"
	"time"
)

/* Function to generate a timestamp */
func generate_timestamp() string {
	//Get the current time
	var currentTime time.Time = time.Now()

	//Merge into string
	var timestamp string = fmt.Sprintf("%s %d, %d at %d:%d:%d",
		currentTime.Month(),
		currentTime.Day(),
		currentTime.Year(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second(),
	)

	//Return
	return "\n<!-- Manifest created by scrawler (" + timestamp + ")-->\n"
}
