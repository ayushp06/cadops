package cad

import "testing"

func TestIsCADExtension(t *testing.T) {
	t.Parallel()

	for _, extension := range []string{".sldprt", ".sldasm", ".prt", ".step", ".stp", ".iges", ".igs", ".stl", ".f3d", ".f3z", ".ipt", ".iam", ".fcstd", ".kicad_pro", ".kicad_pcb", ".kicad_sch", ".kicad_sym", ".kicad_prl", ".kicad_mod", ".kicad_wks", ".kicad_dru"} {
		if !IsCADExtension(extension) {
			t.Fatalf("expected %s to be detected", extension)
		}
	}

	if !IsCADExtension(".FCSTD") {
		t.Fatalf("expected uppercase extension to be detected")
	}

	if IsCADExtension(".txt") {
		t.Fatalf("did not expect .txt to be detected")
	}
}

func TestSupportedLFSExtensionsSkipsKiCadTextFiles(t *testing.T) {
	t.Parallel()

	for _, extension := range SupportedLFSExtensions() {
		if extension == ".kicad_pcb" || extension == ".kicad_sch" || extension == ".kicad_pro" {
			t.Fatalf("did not expect %s to require Git LFS", extension)
		}
	}
}

func TestIsCADPath(t *testing.T) {
	t.Parallel()

	if !IsCADPath("models/gearbox.SLDASM") {
		t.Fatalf("expected path to be detected as CAD")
	}

	if IsCADPath("README") {
		t.Fatalf("did not expect extensionless file to be CAD")
	}
}
