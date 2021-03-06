package initsh

import (
	"os"
	"path/filepath"
	"strings"
)

func WalkThru (importArg string, dir ImportDirector) error {
	ing, root, err := parseImportStr(importArg, dir)
	if (err != nil) {
		return err
	}
	return filterWalkErrs(walkIngest(root, *ing))
}

func WalkPreScanned (importArg string, dir ImportDirector, paths []string) error {
	ing, root, err := parseImportStr(importArg, dir)
	if (err != nil) {
		return err
	}
	return filterWalkErrs(preScanIngest(root, *ing, paths))
}

func filterWalkErrs (err error) error {
	if (err == filepath.SkipDir) {
		return nil
	} else {
		return err
	}
}

func walkIngest (root string, hndl PathIngester) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if (err != nil) {
			return err
		} else if isIgnorable (path) {
			return ignoreWalk(info)
		} else {
			return ingestWalk(path, hndl)
		}
	})
}

func preScanIngest (root string, hndl PathIngester, paths []string) error {
	for _, path := range filterImportRoot(root, paths) {
		err := preScanSelect(path, hndl)
		if (err != nil) {
			return err
		}
	}
	return nil
}

func filterImportRoot (importRoot string, pathUniv []string) []string {
	var paths []string
	for _, path := range pathUniv {
		if (strings.HasPrefix(path, importRoot)) {
			paths = append(paths, path)
		}
	}
	return paths
	
}

func FilterImportDirective (importArg string, pathUniv []string) ([]string, error) {
	lex, err := lexImportStr(importArg)
	root := lex.importPath
	if (err != nil) {
		return make([]string, 0), err
	}
	return filterImportRoot(root, pathUniv), nil
}

func preScanSelect (path string, hndl PathIngester) error {
	if (isIgnorable(path)) {
		return nil
	} else {
		return ingestWalk(path, hndl)
	}
}

func ignoreWalk (info os.FileInfo) error {
	if isSymLink(info.Mode()) {
		return nil
	} else if (info.IsDir()) {
		return filepath.SkipDir
	} else {
		return nil
	}
}

func ingestWalk (path string, hndl PathIngester) error {
	linkInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	return hndl.ingestPath(path, linkInfo)
}

func isIgnorable (path string) bool {
	return isBrokenLink(path) || isIgnorableName(path)
}

func isIgnorableName (path string) bool {
	base := filepath.Base(path)
	return isHiddenBase(base) || isScratchBase(base) || isInitShConfig(base)
}

func isInitShConfig (base string) bool {
	return base == "init.sh"
}

func isHiddenBase (path string) bool {
	return strings.HasPrefix(path, ".") &&
		path != "." && path != "./" &&
		path != ".." && path != "../"
}

func isScratchBase (path string) bool {
	return strings.HasSuffix(path, "~")
}
