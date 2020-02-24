// Copyright (c) 2019 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Configuration holds the agent configuration options
type Configuration struct {
	Hostname               string
	Envtype                string
	Location               string
	Serverurl              string
	Serverusr              string
	Serverpsw              string
	Frequency              int
	HostType               string
	EnableServerValidation bool
	ParallelizeRequests    bool
	Hypervisors            []Hypervisor
}

// Hypervisor holds the parameters used to connect to the hypervisor
type Hypervisor struct {
	Type       string
	Endpoint   string
	Username   string
	Password   string
	OvmUserKey string
	OvmControl string
}

// ReadConfig reads the configuration file from the current dir
// or /opt/ercole-agent
func ReadConfig() Configuration {

	baseDir := GetBaseDir()

	configFile := baseDir + "/config.json"
	ex := exists(configFile)
	if !ex {
		configFile = "/opt/ercole-agent-virtualization/config.json"
	}

	raw, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Fatal("Unable to read configuration file: ", err)
	}

	var conf Configuration
	err = json.Unmarshal(raw, &conf)

	if err != nil {
		log.Fatal("Unable to parse configuration file: ", err)
	}

	return conf

}

func exists(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		return false
	}
	return true
}

// GetBaseDir return path to the directory where is the executable
func GetBaseDir() string {
	s, _ := os.Readlink("/proc/self/exe")
	s = filepath.Dir(s)

	return s
}
