package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	confauth "github.com/zk101/nixie/app/auth/config"
	confchat "github.com/zk101/nixie/app/chat/config"
	confloadtest "github.com/zk101/nixie/app/loadtest/config"
	conftelemetry "github.com/zk101/nixie/app/telemetry/config"
	confws "github.com/zk101/nixie/app/ws/config"
	"github.com/zk101/nixie/lib/buildinfo"
	libconf "github.com/zk101/nixie/lib/config"
	"github.com/zk101/nixie/lib/consul"
)

type testConfig struct {
	TestBool    bool
	TestFloat32 float32
	TestFloat64 float64
	TestInt     int
	TestInt8    int8
	TestInt16   int16
	TestInt32   int32
	TestInt64   int64
	TestUint    uint
	TestUint8   uint8
	TestUint16  uint16
	TestUint32  uint32
	TestUint64  uint64
	TestString  string
}

var (
	configFile       string
	configType       string
	consulClient     *consul.Client
	consulAddress    string
	consulDatacenter string
	consulPrefix     string
	consulScheme     string
	consulToken      string
	debug            bool
)

func init() {
	flag.StringVar(&configFile, "configfile", "", "Config File to process")
	flag.StringVar(&configType, "configtype", "", "Config Type to process")
	flag.StringVar(&consulAddress, "consuladdr", "localhost:8500", "Consul Address")
	flag.StringVar(&consulDatacenter, "consuldc", "global", "Consul Datacentre")
	flag.StringVar(&consulPrefix, "consulprefix", "", "Consul Prefix")
	flag.StringVar(&consulScheme, "consulscheme", "http", "Consul Scheme")
	flag.StringVar(&consulToken, "consultoken", "", "Consul Token")
	flag.BoolVar(&debug, "debug", false, "Debug mode. Show actions, do nothing.")

	var displayver = flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayver == true {
		fmt.Printf("Build Type: %s\n", buildinfo.BuildType)
		fmt.Printf("Build TimeStamp: %s\n", buildinfo.BuildStamp)
		fmt.Printf("Build Revision: %s\n", buildinfo.BuildRevision)
		os.Exit(0)
	}

	if configFile == "" {
		log.Fatalln("Require a config file")
	}

	if configType == "" {
		log.Fatalln("Require a config type")
	}

	consulClient = consul.NewClient(&consul.Config{Address: consulAddress, Scheme: consulScheme, Datacenter: consulDatacenter, Token: consulToken})
}

func main() {
	switch configType {
	case "auth":
		conf := confauth.NewConfig()
		processConfig(conf)
	case "chat":
		conf := confchat.NewConfig()
		processConfig(conf)
	case "loadtest":
		conf := confloadtest.NewConfig()
		processConfig(conf)
	case "telemetry":
		conf := conftelemetry.NewConfig()
		processConfig(conf)
	case "test":
		conf := &testConfig{}
		processConfig(conf)
	case "ws":
		conf := confws.NewConfig()
		processConfig(conf)

	default:
		log.Fatalln("Unsupported Config Type")
	}
}

// processConfig loads the config and sets up the base values and calls processRecurse
func processConfig(conf interface{}) {
	if err := libconf.Load(configFile, conf); err != nil {
		log.Fatalf("The specified Configuration File could not be loaded: %s\n", err.Error())
	}

	reflectConf := reflect.ValueOf(conf).Elem()
	var keyBase string

	if consulPrefix == "" {
		keyBase = reflectConf.Type().Name()
	} else {
		keyBase = consulPrefix + "/" + reflectConf.Type().Name()
	}

	processRecurse(reflectConf, keyBase)
}

// processRecurse is a recursive function to process a struct
func processRecurse(structValue reflect.Value, key string) {
	structToProcess := structValue.Type()

	for i := 0; i < structToProcess.NumField(); i++ {
		switch structToProcess.Field(i).Type.Kind() {
		case reflect.Struct:
			processRecurse(structValue.Field(i), key+"/"+structToProcess.Field(i).Name)
		default:
			processValue(structValue.Field(i), key+"/"+structToProcess.Field(i).Name)
		}
	}
}

// processValue extracts the value and sticks it into consul
func processValue(reflectValue reflect.Value, key string) {
	var value string

	reflectType := reflectValue.Kind()
	switch reflectType {
	case reflect.Bool:
		value = strconv.FormatBool(reflectValue.Interface().(bool))

	case reflect.Float32:
		value = strconv.FormatFloat(float64(reflectValue.Interface().(float32)), 'f', -1, 32)

	case reflect.Float64:
		value = strconv.FormatFloat(reflectValue.Interface().(float64), 'f', -1, 64)

	case reflect.Int:
		value = strconv.FormatInt(int64(reflectValue.Interface().(int)), 10)

	case reflect.Int8:
		value = strconv.FormatInt(int64(reflectValue.Interface().(int8)), 10)

	case reflect.Int16:
		value = strconv.FormatInt(int64(reflectValue.Interface().(int16)), 10)

	case reflect.Int32:
		value = strconv.FormatInt(int64(reflectValue.Interface().(int32)), 10)

	case reflect.Int64:
		value = strconv.FormatInt(reflectValue.Interface().(int64), 10)

	case reflect.Uint:
		value = strconv.FormatUint(uint64(reflectValue.Interface().(uint)), 10)

	case reflect.Uint8:
		value = strconv.FormatUint(uint64(reflectValue.Interface().(uint8)), 10)

	case reflect.Uint16:
		value = strconv.FormatUint(uint64(reflectValue.Interface().(uint16)), 10)

	case reflect.Uint32:
		value = strconv.FormatUint(uint64(reflectValue.Interface().(uint32)), 10)

	case reflect.Uint64:
		value = strconv.FormatUint(reflectValue.Interface().(uint64), 10)

	case reflect.String:
		value = reflectValue.Interface().(string)
	}

	if debug == true {
		fmt.Println(strings.ToLower(key))
		fmt.Println(value)
		return
	}

	if err := consulClient.PutKV(strings.ToLower(key), value); err != nil {
		log.Printf("Consul PutKV Error: %s", err.Error())
	}
}

// EOF
