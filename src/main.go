package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

/* Constants */
const DTD string = "<!DOCTYPE service_bundle SYSTEM '/usr/share/lib/xml/dtd/service_bundle.dtd.1'>"

/* Boolean settings */
const version float32 = 1.80

/* Main */
func main() {

	//Set cmd args up and parse flags

	//Init optional command line args
	var cmd_i_ptr *bool
	var cmd_V_ptr *bool

	//Delegate certain flags to illumos platforms only; else, set false
	if runtime.GOOS == "illumos" {
		cmd_i_ptr = flag.Bool("i", false, "automatically install generated file on exit")
		cmd_V_ptr = flag.Bool("V", false, "validate generated file on exit")
	}

	//Set general command line args
	cmd_d_ptr := flag.String("d", "NODEPS", "colon separated dependencies, with substrings separated by '@' symbols")
	cmd_o_ptr := flag.String("o", "MyService", "specify output file name")
	cmd_q_ptr := flag.Bool("q", false, "suppress output")
	cmd_s_ptr := flag.String("s", "NOARGS", "colon separated S-args, with substrings separated by '@' symbols")
	cmd_v_ptr := flag.Bool("v", false, "print program version and exit")
	cmd_x_ptr := flag.Bool("x", false, "print to stdout")

	//Parse command line args
	flag.Parse()

	/*
		CMD Args:
		-i : automatically install generated file
		-o (string) : specify output
		-q : quiet mode, suppressed output
		-s (name=value) : specify a name / value pair; needs a list of acceptable values pre-made
		-V : validate on exit
		-v : print version and exit
	*/

	//Initialize map for -s name/value pairs
	var s_args map[string]string = make(map[string]string)
	s_args["start-method"] = ":true"
	s_args["stop-method"] = ":true"
	s_args["restart-method"] = ":true"
	s_args["service-name"] = *cmd_o_ptr
	s_args["service-description"] = "MyService Description."
	s_args["timeout-seconds"] = "60"

	//Convert os.Args to list, exclude binary name
	var args map[int]string = make(map[int]string)
	for index, value := range os.Args[1:] {
		args[index] = value
	}

	//Determine boolean args
	var output_file string = fmt.Sprintf("./%s.xml", s_args["service-name"])
	if *cmd_v_ptr {
		fmt.Printf("svcbundle version %.2f %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	//Check for exclusive args
	if cmd_i_ptr != nil && *cmd_i_ptr && *cmd_o_ptr != "MyService.xml" {
		fmt.Println("error: '-i' and '-o' are exclusive arguments!")
		os.Exit(2)
	} else if *cmd_q_ptr && *cmd_x_ptr {
		fmt.Println("error: '-q' and '-x' are exclusive arguments!")
		os.Exit(2)
	}

	//Iterate over and store -s name/value pairs
	if *cmd_s_ptr != "NOARGS" {
		//Check the S-arg string contains a colon; must end in a colon even if one argument is used
		if strings.Contains(*cmd_s_ptr, ":") || (strings.Count(*cmd_s_ptr, "@") == 1) {
			tmpString := strings.Split(*cmd_s_ptr, ":")
			for i := 0; i < len(tmpString); i++ {
				//Check the substring contains an @ symbol
				if strings.Contains(tmpString[i], "@") {
					tmpSubstring := strings.Split(tmpString[i], "@")
					//Check is s-arg is valid
					if _, ok := s_args[tmpSubstring[0]]; !ok {
						fmt.Println("BAD KEY: " + tmpSubstring[1])
						os.Exit(1)
					} else {
						s_args[tmpSubstring[0]] = tmpSubstring[1]
						//Update output file if new service name is provided, will be overwritten if "-o" argument is present
						if tmpSubstring[0] == "service-name" && *cmd_o_ptr != "MyService.xml" {
							output_file = fmt.Sprintf("./%s.xml", s_args["service-name"])
						}
						//Determine if quiet mode is enabled or not
						if !*cmd_q_ptr {
							fmt.Printf("Setting '%s' to '%s'.\n", tmpSubstring[0], tmpSubstring[1])
						}
					}
				} else {
					fmt.Println("Substring separator not found!")
					continue
				}

			}
		} else {
			fmt.Println("S-arg separator not found!")
		}

	}

	//Update filepath if autoinstall is true
	if cmd_i_ptr != nil && *cmd_i_ptr {
		output_file = "/lib/svc/manifest/system/" + output_file
		if !*cmd_q_ptr {
			fmt.Println("Set program to automatically install manifest after completion.")
		}
	}

	//Create service_bundle instance and add a service
	var svcbundle *service_bundle = &service_bundle{Name: s_args["service-name"], Type: "manifest"}
	svcbundle.Service = add_service(s_args)

	//Update dependencies as needed
	if *cmd_d_ptr != "NODEPS" {
		//Check for substring separator
		if strings.Contains(*cmd_d_ptr, "@") {
			//Check for dep separator
			if strings.Contains(*cmd_d_ptr, ":") {
				//Split substrings if they exist
				tmpString := strings.Split(*cmd_d_ptr, ":")
				//Loop over substrings, then split
				for i := 0; i < len(tmpString); i++ {
					tmpSubstring := strings.Split(tmpString[i], "@")
					svcbundle.Service.Dependency[i] = add_dep(tmpSubstring[0], tmpSubstring[1], tmpSubstring[2], tmpSubstring[3], tmpSubstring[4])
				}
			} else {
				//Loop over substrings, then split
				tmpSubstring := strings.Split(*cmd_d_ptr, "@")
				svcbundle.Service.Dependency = []*dependency{add_dep(tmpSubstring[0], tmpSubstring[1], tmpSubstring[2], tmpSubstring[3], "svc:"+tmpSubstring[4])}
			}
		} else {
			//Print error and exit
			fmt.Println("Substring path separator not found.")
			os.Exit(3)
		}

	} else {
		//Use default service setup - just depend on multi user runlevel
		svcbundle.Service.Dependency = []*dependency{add_dep("multi_user_dependency", "require_all", "none", "service", "svc:/milestone/multi-user")}
	}

	//Marshall XML
	out, err := xml.MarshalIndent(svcbundle, " ", "  ")
	err_check(err)
	var tmpString string = strings.Replace(string(out), `<loctext xml:lang="C"></loctext>`, `<loctext xml:lang="C">`+s_args["service-name"]+`</loctext>`, 1)
	tmpString = xml.Header + DTD + generate_timestamp() + strings.Replace(tmpString, `<loctext xml:lang="C"></loctext>`, `<loctext xml:lang="C">`+s_args["service-description"]+`</loctext>`, 1)

	//Open a file for writing
	if !*cmd_x_ptr {
		file, err := os.Create(output_file)
		err_check(err)
		defer file.Close()

		//Add loctext content, write to file, return absolute path
		_, err = file.WriteString(tmpString)
		err_check(err)
		abspath, err := filepath.Abs(output_file)
		err_check(err)
		if !*cmd_q_ptr {
			fmt.Println("Wrote file to : " + abspath)
		}

		//Validate on exit if -V flag supplied
		if cmd_V_ptr != nil && *cmd_V_ptr {
			validate(abspath)
		}
	} else {
		fmt.Println(tmpString)
	}

}

/* Function to check errors */
func err_check(e error) {
	if e != nil {
		panic(e)
	}
}
