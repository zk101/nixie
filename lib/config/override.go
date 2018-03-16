package config

import (
	"reflect"
	"regexp"
	"strconv"
)

// kvstore provides abstraction to a key value store
type kvstore interface {
	Get(string) (string, error)
}

func overrideCheck(reflectValue reflect.Value, key string, kv kvstore) {
	reflectType := reflectValue.Kind()

	overrideValue, err := kv.Get(key)
	if err != nil {
		return
	}

	if overrideValue != "" {
		switch reflectType {
		case reflect.Bool:
			matched, err := regexp.MatchString(`^((?i)true|1)$`, overrideValue)
			if matched == true && err == nil {
				reflectValue.SetBool(true)
			}
			matched, err = regexp.MatchString(`^((?i)false|0)$`, overrideValue)
			if matched == true && err == nil {
				reflectValue.SetBool(false)
			}

		case reflect.Float32, reflect.Float64:
			matched, err := regexp.MatchString(`^-{0,1}[0-9]{1,}\.[0-9]{1,}$`, overrideValue)
			if matched == true && err == nil {
				overrideValueFloat, err := strconv.ParseFloat(overrideValue, 64)
				if err != nil {
					return
				}
				reflectValue.SetFloat(overrideValueFloat)
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			matched, err := regexp.MatchString(`^-{0,1}[0-9]+$`, overrideValue)
			if matched == true && err == nil {
				overrideValueInt, err := strconv.ParseInt(overrideValue, 10, 64)
				if err != nil {
					return
				}
				reflectValue.SetInt(overrideValueInt)
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			matched, err := regexp.MatchString(`^-{0,1}[0-9]+$`, overrideValue)
			if matched == true && err == nil {
				overrideValueInt, err := strconv.ParseUint(overrideValue, 10, 64)
				if err != nil {
					return
				}
				reflectValue.SetUint(overrideValueInt)
			}

		case reflect.String:
			reflectValue.SetString(overrideValue)
		}
	}
}

// overrideRecurse is a recursive function to process a struct
func overrideRecurse(structValue reflect.Value, key, seperator string, kv kvstore) {
	structToProcess := structValue.Type()

	for i := 0; i < structToProcess.NumField(); i++ {
		switch structToProcess.Field(i).Type.Kind() {
		case reflect.Struct:
			overrideRecurse(structValue.Field(i), key+seperator+structToProcess.Field(i).Name, seperator, kv)
		default:
			overrideCheck(structValue.Field(i), key+seperator+structToProcess.Field(i).Name, kv)
		}
	}
}

// Override takes a structure and checks kvstore for a matching key and attempts to replace its value
func Override(config interface{}, prefix, seperator string, kv kvstore) {
	if seperator == "" {
		return
	}

	if config == nil {
		return
	}

	if kv == nil {
		return
	}

	reflectConf := reflect.ValueOf(config).Elem()
	var keyBase string

	if prefix == "" {
		keyBase = reflectConf.Type().Name()
	} else {
		keyBase = prefix + seperator + reflectConf.Type().Name()
	}

	overrideRecurse(reflectConf, keyBase, seperator, kv)
}

// EOF
