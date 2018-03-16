package buildinfo

/*
	This package simply contains vars that are set at build using -X.

	CGO_ENABLED=0 go build -o ${FOLDER_BIN}/${OPT_APP} -ldflags "-w -s -extldflags '-static' -X 'github.com/zk101/nixie/lib/buildinfo.JenkinsBuild=123'"
*/

var (
	// BuildType is what system build the binary
	BuildType string

	// BuildRevision is the mercurial or git revision hash at build time
	BuildRevision string

	// BuildStamp is the DataTime string when the build occured
	BuildStamp string
)

// EOF
