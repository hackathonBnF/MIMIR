/*
* @Author: Bartuccio Antoine
* @Date:   2018-09-16 17:18:40
* @Last Modified by:   klmp200
* @Last Modified time: 2018-10-01 23:44:21
 */
package settings

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

type settings map[string]interface{}

var cfg settings
var cfgLock sync.Mutex

func loadSettings(path string) error {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	return nil
}

func InitSettings(settingsPath, settingsCustomPath string) {

	cfg = make(settings)
	cfgLock = sync.Mutex{}
	cfgLock.Lock()
	defer cfgLock.Unlock()
	if err := loadSettings(settingsPath); err != nil {
		panic(err)
	}
	if err := loadSettings(settingsCustomPath); err != nil {
		log.Println(err)
	}
}

func SettingsValue(key string) interface{} {

	cfgLock.Lock()
	defer cfgLock.Unlock()
	return cfg[key]
}
