package builder

var Languages = map[string]Language{
	"py":   Python,
	"node": Node,
	"ts":   TypeScript,
	"rb":   Ruby,
	"c":    C,
	"cs":   CSharp,
	"java": Java,
	"wasm": WebAssembly,
	"cob":  Cobol,
	"file": File,
	"mock": Mock,
	"rpc":  RPC,
}

type Language int

const (
	Python Language = iota
	Node
	TypeScript
	Ruby
	C
	CSharp
	Java
	WebAssembly
	Cobol
	File
	Mock
	RPC
)
