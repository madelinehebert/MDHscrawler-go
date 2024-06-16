package main

import (
	"errors"
	"fmt"
	"os"
)

const START string = "#! /usr/bin/bash\nnohup SVCBIN &"
const STOP string = "#! /usr/bin/bash\nkill -9 `pgrep SVCBIN`"

func generate_scripts(install_location string, name string) error {
	fmt.Println(install_location)
	//Generate start.sh file
	startFile, err := os.Create(fmt.Sprintf("%sStart_%s.sh", install_location, name))
	if err != nil {
		return errors.New("error opening start file: " + err.Error())
	}
	defer startFile.Close()

	//Write in start.sh boilerplate code
	if _, err := startFile.WriteString(START); err != nil {
		return errors.New("error writing code to start file: " + err.Error())
	}

	//Generate start.sh file
	stopFile, err := os.Create(fmt.Sprintf("%sStop_%s.sh", install_location, name))
	if err != nil {
		return errors.New("error opening stop file: " + err.Error())
	}
	defer stopFile.Close()

	//Write in start.sh boilerplate code
	if _, err := stopFile.WriteString(STOP); err != nil {
		return errors.New("error writing code to stop file: " + err.Error())
	}

	//Is all is well, return nil
	return nil
}
