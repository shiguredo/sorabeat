package cmd

import (
	cmd "github.com/elastic/beats/libbeat/cmd"
	"github.com/elastic/beats/metricbeat/beater"

	// import modules of sorabeat
	_ "github.com/shiguredo/sorabeat/include"
)

// Name of this beat
var Name = "sorabeat"

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmd(Name, "", beater.New)
