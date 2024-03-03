package main

import "encoding/xml"

/* Structs */

// Struct for top level service_bundle
type service_bundle struct {
	XMLName xml.Name `xml:"service_bundle"`

	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Service *service `xml:service`
}

// Struct for services
type service struct {
	XMLName xml.Name `xml:"service"`

	Name           string            `xml:"name,attr"`
	Type           string            `xml:"type,attr"`
	Dependency     []*dependency     `xml:dependency`
	Exec_Method    []*exec_method    `xml:"exec_method"`
	Property_Group []*property_group `xml:"property_group"`
	Instance       *instance         `xml:"instance"`
	Template       *template         `xml:"template"`
	Version        float32           `xml:"version,attr"`
}

// Struct for dependencies
type dependency struct {
	XMLName xml.Name `xml:"dependency"`

	Name         string        `xml:"name,attr"`
	Type         string        `xml:"type,attr"`
	Grouping     string        `xml:"grouping,attr"`
	Restart_On   string        `xml:"restart_on,attr"`
	Service_FMRI *service_fmri `xml:service_fmri`
}

// Struct for service_fmris
type service_fmri struct {
	XMLName xml.Name `xml:"service_fmri"`

	Value string `xml:"value,attr"`
}

// Struct for exec_methods
type exec_method struct {
	XMLName xml.Name `xml:exec_method`

	Name            string `xml:"name,attr"`
	Type            string `xml:"type,attr"`
	Exec            string `xml:"exec,attr"`
	Timeout_Seconds string `xml:"timeout_seconds,attr"`
}

// Struct for property_groups
type property_group struct {
	XMLName xml.Name `xml:"property_group"`

	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	PropVal *propval `xml:propval`
}

// Struct for propvals
type propval struct {
	XMLName xml.Name `xml:"propval"`

	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
}

// Struct for instances
type instance struct {
	XMLName xml.Name `xml:"instance"`

	Name    string `xml:"name,attr"`
	Enabled bool   `xml:"enabled,attr"`
}

// Struct for templates
type template struct {
	XMLName xml.Name `xml:"template"`

	Common_Name *common_name `xml:"common_name"`
	Description *description `xml:"description"`
}

// Struct for common_names
type common_name struct {
	XMLName xml.Name `xml:"common_name"`

	LocText *loctext `xml:"loctext"`
}

// Struct for descriptions
type description struct {
	XMLName xml.Name `xml:"description"`

	LocText *loctext `xml:"loctext"`
}

// Struct for loctexts
type loctext struct {
	XMLName xml.Name `xml:"loctext"`

	XMLLang string `xml:"xml:lang,attr"`
}
