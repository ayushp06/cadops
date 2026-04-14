package cad

// FileType describes a supported CAD file extension.
type FileType struct {
	Extension        string
	Name             string
	UseLFS           bool
	RecommendLocking bool
}
