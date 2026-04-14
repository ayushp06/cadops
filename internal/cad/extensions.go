package cad

import "strings"

var supportedTypes = []FileType{
	{Extension: ".sldprt", Name: "SolidWorks Part", UseLFS: true, RecommendLocking: true},
	{Extension: ".sldasm", Name: "SolidWorks Assembly", UseLFS: true, RecommendLocking: true},
	{Extension: ".step", Name: "STEP", UseLFS: true, RecommendLocking: false},
	{Extension: ".stp", Name: "STEP", UseLFS: true, RecommendLocking: false},
	{Extension: ".iges", Name: "IGES", UseLFS: true, RecommendLocking: false},
	{Extension: ".igs", Name: "IGES", UseLFS: true, RecommendLocking: false},
	{Extension: ".ipt", Name: "Inventor Part", UseLFS: true, RecommendLocking: true},
	{Extension: ".iam", Name: "Inventor Assembly", UseLFS: true, RecommendLocking: true},
	{Extension: ".fcstd", Name: "FreeCAD Document", UseLFS: true, RecommendLocking: true},
}

var byExtension = func() map[string]FileType {
	index := make(map[string]FileType, len(supportedTypes))
	for _, fileType := range supportedTypes {
		index[fileType.Extension] = fileType
	}
	return index
}()

// SupportedTypes returns the built-in CAD registry.
func SupportedTypes() []FileType {
	out := make([]FileType, len(supportedTypes))
	copy(out, supportedTypes)
	return out
}

// SupportedExtensions returns the configured CAD extensions in stable order.
func SupportedExtensions() []string {
	out := make([]string, 0, len(supportedTypes))
	for _, fileType := range supportedTypes {
		out = append(out, fileType.Extension)
	}
	return out
}

// Lookup returns the CAD file type for an extension.
func Lookup(extension string) (FileType, bool) {
	fileType, ok := byExtension[strings.ToLower(extension)]
	return fileType, ok
}

// IsCADExtension reports whether the extension belongs to the CAD registry.
func IsCADExtension(extension string) bool {
	_, ok := Lookup(extension)
	return ok
}

// IsCADPath reports whether the file path has a supported CAD extension.
func IsCADPath(path string) bool {
	lastDot := strings.LastIndex(path, ".")
	if lastDot == -1 {
		return false
	}
	return IsCADExtension(path[lastDot:])
}
