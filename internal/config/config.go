package config

import (
	"errors"
	"reflect"
	"strings"
)

// stringToMapstringHookFunc will decode a string to a mapstring.
func stringToMapstringHookFunc(f, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String || t != reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf("")) {
		return data, nil
	}

	mapstring := make(map[string]string)
	for _, elements := range strings.Split(data.(string), ",") {
		config := strings.Split(elements, ":")
		if len(config) != 2 {
			return nil, errors.New("expected string of format 'key0:value0,key1:value1,...,keyN:valueN'")
		}
		mapstring[config[0]] = config[1]
	}

	return mapstring, nil
}
