package cad

import "testing"

func TestIsCADExtension(t *testing.T) {
	t.Parallel()

	if !IsCADExtension(".sldprt") {
		t.Fatalf("expected .sldprt to be detected")
	}

	if !IsCADExtension(".FCSTD") {
		t.Fatalf("expected uppercase extension to be detected")
	}

	if IsCADExtension(".txt") {
		t.Fatalf("did not expect .txt to be detected")
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
