package main

// Create and return a service
func add_service(s_args map[string]string) *service {
	//Create service child and assign to service_bundle parent
	var svc *service = &service{Name: s_args["service-name"], Type: "service", Version: 1.0}

	//Create dependency for service
	var dep *dependency = &dependency{Name: "multi_user_dependency", Grouping: "require_all", Restart_On: "none", Type: "service"}
	var svcFMRI *service_fmri = &service_fmri{Value: "svc:/milestone/multi-user"}
	dep.Service_FMRI = svcFMRI
	svc.Dependency = []*dependency{dep}

	//Create start, stop, and restart exec methods
	var exec_start *exec_method = &exec_method{Name: "start", Exec: s_args["start-method"], Timeout_Seconds: s_args["timeout-seconds"], Type: "method"}
	var exec_stop *exec_method = &exec_method{Name: "stop", Exec: s_args["stop-method"], Timeout_Seconds: s_args["timeout-seconds"], Type: "method"}
	var exec_restart *exec_method = &exec_method{Name: "restart", Exec: s_args["restart-method"], Timeout_Seconds: s_args["timeout-seconds"], Type: "method"}
	svc.Exec_Method = []*exec_method{exec_start, exec_stop, exec_restart}

	//Create property_group for service
	var propgroup *property_group = &property_group{Name: "startd", Type: "framework"}
	var propval *propval = &propval{Name: "duration", Type: "astring", Value: "transient"}
	propgroup.PropVal = propval
	svc.Property_Group = []*property_group{propgroup}

	//Add an instance
	var inst *instance = &instance{Name: "default", Enabled: true}
	svc.Instance = inst

	//Add a template
	var template *template = &template{}
	var commonname *common_name = &common_name{LocText: &loctext{XMLLang: "C"}}
	var desc *description = &description{LocText: &loctext{XMLLang: "C"}}
	template.Common_Name = commonname
	template.Description = desc
	svc.Template = template

	//return svc to caller function
	return svc
}
