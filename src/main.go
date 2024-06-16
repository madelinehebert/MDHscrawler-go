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
const version float32 = 1.90

/* Other data */
var INSTALL_DIR string = "./"

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
	cmd_o_ptr := flag.String("o", "MyService.xml", "specify output file name; file extention is automatically added")
	cmd_q_ptr := flag.Bool("q", false, "suppress output")
	cmd_s_ptr := flag.String("s", "NOARGS", "colon separated S-args, with substrings separated by '@' symbols")
	cmd_S_ptr := flag.Bool("S", false, "generate start and stop scripts for the service, and update those methods")
	cmd_v_ptr := flag.Bool("v", false, "print program version and exit")
	cmd_x_ptr := flag.Bool("x", false, "print to stdout")
	cmd_help_ptr := flag.Bool("?", false, "print the help menu, version, and exit")

	//Parse command line args
	flag.Parse()

	//Print help menu if needed
	if *cmd_help_ptr {
		//Print version header
		printv()
		//Help menu goes here
		flag.PrintDefaults()
		os.Exit(0)
	}

	//Check version arg
	if *cmd_v_ptr {
		printv()
		os.Exit(0)
	}

	//Initialize map for -s name/value pairs
	var s_args map[string]string = make(map[string]string)
	s_args["start-method"] = ":true"
	s_args["stop-method"] = ":true"
	s_args["restart-method"] = ":true"
	s_args["service-name"] = "MyService"
	s_args["service-description"] = "MyService Description."
	s_args["timeout-seconds"] = "60"

	//Convert os.Args to list, exclude binary name
	var args map[int]string = make(map[int]string)
	for index, value := range os.Args[1:] {
		args[index] = value
	}

	//Set output file name
	var output_file string = *cmd_o_ptr

	//Check for exclusive args
	if *cmd_q_ptr && *cmd_x_ptr {
		fmt.Println("error: '-q' and '-x' are exclusive arguments!")
		os.Exit(2)
	} else if cmd_V_ptr != nil && *cmd_V_ptr && *cmd_x_ptr {
		fmt.Println("error: '-V' and '-x' are exclusive arguments!")
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
						//Check if -S is set, meaning ignore manually supplied start and stop methods
						if *cmd_S_ptr && (tmpSubstring[0] == "start-method" || tmpSubstring[0] == "stop-method") {
							if !*cmd_q_ptr {
								fmt.Println("Ignoring Start and Stop methods; -S is set to true.")
							}
							continue
						}
						s_args[tmpSubstring[0]] = tmpSubstring[1]
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
		INSTALL_DIR = "/lib/svc/manifest/system/"
		output_file = INSTALL_DIR + output_file
		if !*cmd_q_ptr {
			fmt.Println("Set program to automatically install manifest after completion.")
		}
	}

	//Check if scripts should be generated
	if *cmd_S_ptr {
		//Generate scripts
		if err := generate_scripts(INSTALL_DIR, s_args["service-name"]); err != nil {
			fmt.Println(err)
			os.Exit(5)
		}
		//Update start and stop methods
		s_args["start-method"] = fmt.Sprintf("%sStart_%s.sh", INSTALL_DIR, s_args["service-name"])
		s_args["stop-method"] = fmt.Sprintf("%sStop_%s.sh", INSTALL_DIR, s_args["service-name"])
	}

	//Create service_bundle instance and add a service
	var scrawler *service_bundle = &service_bundle{Name: s_args["service-name"], Type: "manifest"}
	scrawler.Service = add_service(s_args)

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
					scrawler.Service.Dependency[i] = add_dep(tmpSubstring[0], tmpSubstring[1], tmpSubstring[2], tmpSubstring[3], tmpSubstring[4])
				}
			} else {
				//Make an array out of the lone pair provided
				tmpSubstring := strings.Split(*cmd_d_ptr, "@")
				scrawler.Service.Dependency = []*dependency{add_dep(tmpSubstring[0], tmpSubstring[1], tmpSubstring[2], tmpSubstring[3], "svc:"+tmpSubstring[4])}
			}
		} else {
			//Print error and exit
			fmt.Println("Substring path separator not found.")
			os.Exit(3)
		}
	} else {
		//Use default service setup - just depend on multi user runlevel
		scrawler.Service.Dependency = []*dependency{add_dep("multi_user_dependency", "require_all", "none", "service", "svc:/milestone/multi-user")}
	}

	//Marshall XML
	out, err := xml.MarshalIndent(scrawler, " ", "  ")
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
			if !*cmd_q_ptr {
				if validate(abspath) {
					fmt.Println("Validated service manifest!")
				} else {
					fmt.Println("Failed to validate service manifest!")
					os.Exit(4)
				}
			} else {
				if !validate(abspath) {
					fmt.Println("Failed to validate service manifest!")
					os.Exit(4)
				}
			}
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
