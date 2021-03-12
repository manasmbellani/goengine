package main

// Target to execute checks against
type Target struct {
	Host     string
	Protocol string
	Port     string
	Path     string
	Basepath string
	Querystr string
	Target   string
}

// CheckToExec is the method, check to execute
type CheckToExec struct {
	CheckID  string
	MethodID string
	Target   Target
	Method   MethodStruct
}

// CheckStruct is a Single check from the Checksfile
type CheckStruct struct {
	ID      string         `yaml:"id"`
	Methods []MethodStruct `yaml:"methods"`
}

// MethodStruct is the method to deploy in method
type MethodStruct struct {
	ID         string   `yaml:"id"`
	Type       string   `yaml:"type"`
	Cmds       []string `yaml:"cmd"`
	CmdDir     string   `yaml:"cmddir"`
	HTTPMethod string   `yaml:"method"`
	Urls       []string `yaml:"urls"`
	Body       []struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	} `yaml:"body"`
	//BodyStr string `yaml:"bodystr"`
	Headers []struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	} `yaml:"headers"`
	Outfile string `yaml:"outfile"`
	Regex   string `yaml:"regex"`
	Notes   string `yaml:"notes"`
}

// ChecksFileStruct defines the structure of the Checks file (in YAML)
type ChecksFileStruct struct {
	Checks []CheckStruct `yaml:"checks"`
}
