package main

/* Function to build an return a dependency struct */
func add_dep(name string, grouping string, restart_on string, svc_type string, fmri string) *dependency {
	/*
		Example setup:
				var dep *dependency = &dependency{Name: "multi_user_dependency", Grouping: "require_all", Restart_On: "none", Type: "service"}
				var svcFMRI *service_fmri = &service_fmri{Value: "svc:/milestone/multi-user"}
	*/

	var dep *dependency = &dependency{Name: name, Grouping: grouping, Restart_On: restart_on, Type: svc_type}
	var svcFMRI *service_fmri = &service_fmri{Value: fmri}
	dep.Service_FMRI = svcFMRI

	return dep
}
