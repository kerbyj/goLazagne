package filesystem

import (
	"os"
	"path/filepath"
	"strings"
)

func getdrives() (r []string) {
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		_, err := os.Open(string(drive) + ":\\")
		if err == nil {
			r = append(r, string(drive))
		}
	}
	return
}

func FindFiles(suffixes []string) []string {

	var (
		interestingFilesList []string

		drives = getdrives()
	)

	for driveNum := range drives {
		var root = drives[driveNum] + ":\\\\"

		var _ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			for i := range suffixes {
				if strings.HasSuffix(info.Name(), suffixes[i]) {
					interestingFilesList = append(interestingFilesList, path)
				}
			}
			return nil
		})
	}

	return interestingFilesList
}
