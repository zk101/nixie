package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/zk101/nixie/lib/buildinfo"
)

// flags holds common cli arguments
type flags struct {
	configFile       string
	consulAddress    string
	consulDatacenter string
	consulPrefix     string
	consulScheme     string
	consulToken      string
	displayVer       bool
}

// processFlags parses flag options and returns a pointer to flags struct
func processFlags() *flags {
	flags := flags{}

	flag.StringVar(&flags.configFile, "configfile", "", "Config File to process")
	flag.StringVar(&flags.consulAddress, "consuladdr", "", "Consul Address")
	flag.StringVar(&flags.consulDatacenter, "consuldc", "global", "Consul Datacentre")
	flag.StringVar(&flags.consulPrefix, "consulprefix", "", "Consul Prefix")
	flag.StringVar(&flags.consulScheme, "consulscheme", "http", "Consul Scheme")
	flag.StringVar(&flags.consulToken, "consultoken", "", "Consul Token")
	flag.BoolVar(&flags.displayVer, "version", false, "Display version and exit")
	flag.Parse()

	if flags.displayVer == true {
		fmt.Printf("Build Type: %s\n", buildinfo.BuildType)
		fmt.Printf("Build TimeStamp: %s\n", buildinfo.BuildStamp)
		fmt.Printf("Build Revision: %s\n", buildinfo.BuildRevision)
		os.Exit(0)
	}

	return &flags
}

// EOF
