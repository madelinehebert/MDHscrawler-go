package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

/* Function to validate the manifest */
func validate(abspath string) bool {
	//Form command string and attach buffer to stdout
	var cmd *exec.Cmd = exec.Command("/usr/sbin/svccfg", "validate", abspath)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run command and check output
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return false
	} else if len(strings.Trim(out.String(), "\n")) == 0 {
		//fmt.Println("Validated!")
		return true
	}

	return false
}
