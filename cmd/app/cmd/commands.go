package cmd

import "github.com/spf13/cobra"

const suppressErrorReporting = "suppress-error-reporting"

func ShouldSuppressErrorReporting(c *cobra.Command) bool {
	for c != nil {
		if _, found := c.Annotations[suppressErrorReporting]; found {
			return true
		}
		c = c.Parent()
	}
	return false
}
