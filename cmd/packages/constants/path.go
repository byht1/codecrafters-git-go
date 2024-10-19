package constants

import "path"

const (
	ROOT_DIR string = ".git"

	PARAM_DEFAULT_VALUE string = ""
)

var (
	ObjectDir string = path.Join(ROOT_DIR, "objects")
	RefsDir   string = path.Join(ROOT_DIR, "refs")
)

func RootPathBuilder(args ...string) string {
	allArgs := append([]string{ROOT_DIR}, args...)
	return path.Join(allArgs...)
}

func ObjectPathBuilder(args ...string) string {
	allArgs := append([]string{ObjectDir}, args...)
	return path.Join(allArgs...)
}

func RefsPathBuilder(args ...string) string {
	allArgs := append([]string{ObjectDir}, args...)
	return path.Join(allArgs...)
}
