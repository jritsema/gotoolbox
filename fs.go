package gotoolbox

import "os"

//IsDirectory returns true if a path is a directory
func IsDirectory(path string) (bool, error) {
	fd, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	switch mode := fd.Mode(); {
	case mode.IsDir():
		return true, nil
	case mode.IsRegular():
		return false, nil
	}
	return false, nil
}

//FileExists returns true if a file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
