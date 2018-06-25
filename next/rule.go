package next

// NewPHPRule create a rule for php
func NewPHPRule(rootDir string, addresses []string) Rule {
	return Rule{
		Root:       rootDir,
		Path:       "/",
		balancer:   &roundRobin{addresses: addresses, index: -1},
		Ext:        ".php",
		SplitPath:  ".php",
		IndexFiles: []string{"index.php"},
	}
}
