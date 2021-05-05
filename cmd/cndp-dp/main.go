/*
 Copyright(c) 2021 Intel Corporation.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultConfigFile = "./config.json"
	devicePrefix      = "cndp"
)

type devicePlugin struct {
	pools map[string]PoolManager
}

type poolConfig struct {
	Name    string   `json:"name"`
	Devices []string `json:"devices"`
}

type config struct {
	Pools []poolConfig `json:"pools"`
}

func main() {
	var configFile string

	flag.Lookup("logtostderr").Value.Set("true")
	flag.StringVar(&configFile, "config", defaultConfigFile, "Location of the device plugin configuration file")
	flag.Parse()

	dp := devicePlugin{
		pools: make(map[string]PoolManager),
	}

	glog.Info("Reading config file " + configFile)
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		glog.Error("Error reading config file " + configFile)
	}

	cfg, err := getConfig(raw)
	if err != nil {
		glog.Error("Error parsing config file " + configFile)
		glog.Fatal(err)
	}

	for _, poolConfig := range cfg.Pools {

		pm := PoolManager{
			Name:          poolConfig.Name,
			Devices:       make(map[string]*pluginapi.Device),
			DpAPISocket:   pluginapi.DevicePluginPath + devicePrefix + "-" + poolConfig.Name + ".sock",
			DpAPIEndpoint: devicePrefix + "-" + poolConfig.Name + ".sock",
			UpdateSignal:  make(chan bool),
		}

		err = pm.Init(poolConfig)
		if err != nil {
			glog.Error("Error initializing pool: " + pm.Name)
			glog.Fatal(err)
		}

		dp.pools[poolConfig.Name] = pm
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case s := <-sigs:
		glog.Infof("Received signal \"%v\"", s)
		for _, pm := range dp.pools {
			glog.Infof("Terminating " + pm.Name)
			pm.Terminate()
		}
		return
	}
}

func getConfig(raw []byte) (*config, error) {
	cfg := &config{}

	err := json.Unmarshal(raw, &cfg)
	if err != nil {
		return nil, err
	}

	glog.Info("Config: " + fmt.Sprintf("%+v\n", cfg))

	return cfg, nil
}
