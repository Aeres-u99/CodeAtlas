package internal

type Symbol struct {
	Name   string `json:"n"`
	Type   string `json:"t"`
	Line   int    `json:"l"`
	Public bool   `json:"p,omitempty"`
	Sig    string `json:"sig,omitempty"`
}

type FileInfo struct {
	Lang    string   `json:"lang"`
	LOC     int      `json:"loc"`
	Imports []string `json:"imports"`
	Symbols []Symbol `json:"symbols"`
}

type Location struct {
	File string `json:"f"`
	Line int    `json:"l"`
}

type Output struct {
	Version   int                 `json:"v"`
	Generated string              `json:"generated"`
	Files     map[string]FileInfo `json:"files"`
	Index     map[string]Location `json:"idx"`
}

type CTag struct {
	Type      string `json:"_type"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Pattern   string `json:"pattern"`
	Line      int    `json:"line"`
	Kind      string `json:"kind"`
	Access    string `json:"access,omitempty"`
	Scope     string `json:"scope,omitempty"`
	ScopeKind string `json:"scopeKind,omitempty"`
}

type AnalysisResult struct {
	Files map[string]FileInfo
	Index map[string]Location
}

type FileAnalysis struct {
	FileInfo FileInfo
	Index    map[string]Location
}
