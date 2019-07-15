package opsmanagerconfig

import (
	"bufio"
	"bytes"
	"github.com/MihaiBojin/terraform-provider-mongodb/mongodb/util"
	"log"
	"strings"

	"github.com/magiconair/properties"
)

// PropertiesFile wrapper for Ops Manager config files
type PropertiesFile struct {
	props *properties.Properties
}

// NewPropertiesFile create a new wrapper for Ops Manager config files
func NewPropertiesFile(data string) *PropertiesFile {
	loader := &properties.Loader{DisableExpansion: true, Encoding: properties.UTF8}
	p, err := loader.LoadBytes([]byte(data))
	util.PanicOnNonNilErr(err)

	log.Print("[DEBUG] Loaded properties file...")
	return &PropertiesFile{props: p}
}

// SetPropertyValue sets/updates a property key, value pair
func (cfg *PropertiesFile) SetPropertyValue(key string, value string) {
	if _, _, err := cfg.props.Set(key, value); err != nil {
		log.Printf("[DEBUG] Could not set property %s=%s; err=%v", key, value, err)
	}
}

// SetComments sets comment(s) for the specified key
func (cfg *PropertiesFile) SetComments(key string, comments []string) {
	cfg.props.SetComments(key, comments)
}

// Write return all properties as a string
func (cfg *PropertiesFile) Write() (string, error) {
	// write all properties to a buffer
	var buf bytes.Buffer
	_, _ = cfg.props.WriteComment(&buf, "#", properties.UTF8)

	// read all lines and fix unnecessary spacing
	// this is a workaround for github.com/magiconair/properties's fixed output format
	reader := bufio.NewReader(bytes.NewReader(buf.Bytes()))
	scanner := bufio.NewScanner(reader)
	var out bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()

		// fix extra spacing before and after "="
		line = strings.Replace(line, " = ", "=", 1)

		out.WriteString(line + "\n")
	}

	return out.String(), nil
}
