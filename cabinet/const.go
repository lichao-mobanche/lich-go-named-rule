package cabinet

// Status type
type Status int

//Status enum
const (
	Complete   Status = iota //All subtags exist
	Incomplete               // All subtags exist but some subtags are Unconnected
	Defective                // Some child tag does not exist
)

// Result TODO
type Result map[string]interface{}

// GroupTag is a special tageName, which use to storage group rule
const GroupTag = "_GroupTag_"
