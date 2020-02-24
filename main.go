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

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/ercole-io/ercole-agent-virtualization/builder"
	"github.com/ercole-io/ercole-agent-virtualization/config"
	"github.com/ercole-io/ercole-agent-virtualization/model"
	"github.com/ercole-io/ercole-agent-virtualization/scheduler"
	"github.com/ercole-io/ercole-agent-virtualization/scheduler/storage"
)

var version = "latest"

func main() {
	rand.Seed(243243)
	configuration := config.ReadConfig()

	doBuildAndSend(configuration)

	memStorage := storage.NewMemoryStorage()
	scheduler := scheduler.New(memStorage)

	_, err := scheduler.RunEvery(time.Duration(configuration.Frequency)*time.Hour, doBuildAndSend, configuration, version)
	if err != nil {
		log.Fatal("Error sending data: ", err)
	}

	scheduler.Start()
	scheduler.Wait()
}

func doBuildAndSend(configuration config.Configuration) {
	hostData := builder.BuildData(configuration, version)
	sendData(hostData, configuration)
}

func sendData(data *model.HostData, configuration config.Configuration) {
	log.Println("Sending data...")

	b, _ := json.Marshal(data)
	s := string(b)

	log.Println("Data:", s)

	hostType := configuration.HostType
	if hostType == "" {
		hostType = "non-defined"
	}

	client := &http.Client{}

	//Disable certificate validation if enableServerValidation is false
	if configuration.EnableServerValidation == false {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	url, err := url.Parse(configuration.Serverurl)
	if err != nil {
		log.Fatal(err)
	} else {
		var ok bool
		_, ok = url.Query()["HostType"]
		if !ok {
			query := url.Query()
			query.Add("HostType", hostType)
			url.RawQuery = query.Encode()
		}
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewReader(b))

	req.Header.Add("Content-Type", "application/json")
	// auth := configuration.Serverusr + ":" + configuration.Serverpsw
	// authEnc := b64.StdEncoding.EncodeToString([]byte(auth))
	// req.Header.Add("Authorization", "Basic "+authEnc)
	req.SetBasicAuth(configuration.Serverusr, configuration.Serverpsw)
	resp, err := client.Do(req)

	sendResult := "FAILED"

	if err != nil {
		log.Println("Error sending data: ", err)
	} else {
		log.Println("Response status: ", resp.Status)
		if resp.StatusCode == 200 {
			sendResult = "SUCCESS"
		}
		defer resp.Body.Close()
	}

	log.Println("Sending result: ", sendResult)
}
