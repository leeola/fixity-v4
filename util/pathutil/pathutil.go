package pathutil

import (
	"fmt"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

func ExpandJoin(paths ...string) (string, error) {
	for i, p := range paths {
		if p != "" {
			ep, err := homedir.Expand(p)
			if err != nil {
				return "", fmt.Errorf("expand: %v", err)
			}
			paths[i] = ep
		}
	}

	return filepath.Join(paths...), nil
}
