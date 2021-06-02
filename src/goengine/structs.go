package main

// Target to execute checks against
type Target struct {
	Outfolder  string
	Host       string
	Protocol   string
	Port       string
	Path       string
	Basepath   string
	Querystr   string
	Target     string
	Folder     string
	AWSProfile string
	AWSRegion  string
	GCPAccount string
	GCPProject string
	GCPRegion  string
	GCPZone    string
	// Useful for lowhanging search
	Company string
}

// CheckToExec is the method, check to execute
type CheckToExec struct {
	CheckID  string
	Target   Target
	CheckDetails  CheckStruct
}

// MethodStruct is the method to deploy in method
type CheckStruct struct {
	Type       string   `yaml:"type"`
	Cmds       []string `yaml:"cmd"`
	CmdDir     string   `yaml:"cmddir"`
	HTTPMethod string   `yaml:"method"`
	Urls       []string `yaml:"url"`
	Body       []struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	} `yaml:"body"`
	//BodyStr string `yaml:"bodystr"`
	Headers []struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	} `yaml:"headers"`
	Keywords       []string `yaml:"keyword"`
	Files          []string `yaml:"file"`
	Outfile        string   `yaml:"outfile"`
	Searches       []string `yaml:"search"`
	WriteToOutfile bool     `yaml:"writetofile"`
	Regex          string   `yaml:"regex"`
	NoRegex        string   `yaml:"noregex"`
	AlertOnMissing bool     `yaml:"alertonmissing"`
	Notes          string   `yaml:"notes"`
}

// ChecksFileStruct defines the structure of the Checks file (in YAML)
type CheckFileStruct struct {
	Check CheckStruct `yaml:"check"`
}
