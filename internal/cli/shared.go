package cli

import (
	"fmt"

	"github.com/cadops/cadops/internal/collab"
)

func printWarnings(warnings []collab.Warning) {
	for _, warning := range warnings {
		fmt.Printf("Warning: %s: %s\n", warning.Title, warning.Details)
	}
}
