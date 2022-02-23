package cmdutils

import "github.com/urfave/cli/v2"

// MergeFlags merges multiple flags group into one
func MergeFlags(flagsGroups ...[]cli.Flag) []cli.Flag {
	allFlags := make([]cli.Flag, 0)
	for _, flagsGroup := range flagsGroups {
		allFlags = append(allFlags, flagsGroup...)
	}

	return allFlags
}
