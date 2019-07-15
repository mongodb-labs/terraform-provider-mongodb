package types

import (
	"fmt"
	"github.com/MihaiBojin/terraform-provider-mongodb/mongodb/util"
	"github.com/hashicorp/terraform/helper/schema"
)

// StrictUnionOrFatal convenience method which throws a fatal error if duplicate keys are found while merging multiple terraform schema maps
func StrictUnionOrFatal(maps ...map[string]*schema.Schema) map[string]*schema.Schema {
	union, err := StrictUnion(maps...)
	util.PanicOnNonNilErr(err)
	return union
}

// StrictUnion merges multiple terraform schema maps, ensuring keys are not repeated
func StrictUnion(maps ...map[string]*schema.Schema) (map[string]*schema.Schema, error) {
	data := make(map[string]*schema.Schema)

	for _, m := range maps {
		if err := mapUnion(&data, m); err != nil {
			// stop here if a duplicate key is found
			return data, err
		}
	}
	return data, nil
}

// mapUnion merges keys from the second parameter into the first map
func mapUnion(collector *map[string]*schema.Schema, data map[string]*schema.Schema) error {
	for k, v := range data {
		// ensure we don't override existing keys
		if _, ok := (*collector)[k]; ok {
			return fmt.Errorf("Duplicate key: %s", k)
		}

		(*collector)[k] = v
	}

	return nil
}

// NewSchemaMap build a new schema map based on the passed configuration
func NewSchemaMap(config ...func() map[string]*schema.Schema) map[string]*schema.Schema {
	data := make(map[string]*schema.Schema)
	for _, configPart := range config {
		err := mapUnion(&data, configPart())
		util.PanicOnNonNilErr(err)
	}

	return data
}

// ReadString reads a string value from an input map by key
func ReadString(input map[string]interface{}, key string) (string, bool) {
	if v, ok := input[key]; ok {
		return v.(string), true
	}
	return "", false
}

// ReadInt reads an int value from an input map by key
func ReadInt(input map[string]interface{}, key string) (int, bool) {
	if v, ok := input[key]; ok {
		return v.(int), true
	}
	return 0, false
}

// ReadFloat reads a float value from an input map by key
func ReadFloat(input map[string]interface{}, key string) (float64, bool) {
	if v, ok := input[key]; ok {
		return v.(float64), true
	}
	return 0.0, false
}

// ReadBool reads a string value from an input map by key
func ReadBool(input map[string]interface{}, key string) (bool, bool) {
	if v, ok := input[key]; ok {
		return v.(bool), true
	}
	return false, false
}

// ReadStringMap reads all values in a string-string map
func ReadStringMap(input map[string]interface{}, key string) (map[string]interface{}, bool) {
	if v, ok := input[key]; ok {
		return v.(map[string]interface{}), true
	}
	return make(map[string]interface{}), false
}
