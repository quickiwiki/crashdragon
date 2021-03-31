// Package config provides the config for CrashDragon
package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

//NilDate stores the string a empty Time type gets printed in Unix mode
const NilDate = "Mon Jan  1 00:00:00 UTC 0001"

//Config holds the structure of the configuration
type Config struct {
	DatabaseConnection string
	UseSocket          bool
	BindAddress        string
	BindSocket         string
	ContentDirectory   string
	TemplatesDirectory string
	AssetsDirectory    string
	SymbolicatorPath   string
	TrimModuleNames    bool
}

//C is the actual configuration read from the file
var C Config

//LoadConfig reads the configuration file in the specified path
func LoadConfig(path string) error {
	_, e := toml.DecodeFile(path, &C)
	return e
}

//WriteConfig writes the current configuration to the specified file
func WriteConfig(path string) error {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(C)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(path), 0750)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, buf.Bytes(), 0600)
	return err
}

//GetConfig loads default values and overwrites them by the ones in a file, or creates a file with them if there is no file
func GetConfig(path string) error {
	//Set default values if there is no config
	C.DatabaseConnection = "host=localhost user=crashdragon dbname=crashdragon sslmode=disable"
	C.UseSocket = false
	C.BindAddress = "0.0.0.0:8080"
	C.BindSocket = "/var/run/crashdragon/crashdragon.sock"
	C.ContentDirectory = "../share/crashdragon/files"
	C.TemplatesDirectory = "./web/templates"
	C.AssetsDirectory = "./web/assets"
	C.SymbolicatorPath = "./minidump_stackwalk"
	C.TrimModuleNames = true

	var cerr error
	if _, err := os.Stat(path); err == nil {
		cerr = LoadConfig(path)
	} else {
		cerr = WriteConfig(path)
	}
	return cerr
}
