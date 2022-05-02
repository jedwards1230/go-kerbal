package dirfs

import (
	"io/fs"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
)

var SkipDir error = fs.SkipDir

type WalkFunc func(path string, info fs.FileInfo, err error) error
type statDirEntry struct {
	info fs.FileInfo
}

func (d *statDirEntry) Name() string               { return d.info.Name() }
func (d *statDirEntry) IsDir() bool                { return d.info.IsDir() }
func (d *statDirEntry) Type() fs.FileMode          { return d.info.Mode().Type() }
func (d *statDirEntry) Info() (fs.FileInfo, error) { return d.info, nil }

// Modified WalkDir function to work with billy.Filesytem instead of os
func WalkDir(repo billy.Filesystem, root string, fn fs.WalkDirFunc) error {
	info, err := repo.Lstat(root)
	if err != nil {
		err = fn(root, nil, err)
	} else {
		err = walkDir(repo, root, &statDirEntry{info}, fn)
	}
	if err == SkipDir {
		return nil
	}
	return err
}

// walk recursively descends path, calling walkFn.
func walkDir(repo billy.Filesystem, path string, d fs.DirEntry, walkDirFn fs.WalkDirFunc) error {
	if err := walkDirFn(path, d, nil); err != nil || !d.IsDir() {
		if err == SkipDir && d.IsDir() {
			// Successfully skipped directory.
			err = nil
		}
		return err
	}

	dirs, err := repo.ReadDir(path)
	if err != nil {
		// Second call, to report ReadDir error.
		err = walkDirFn(path, d, err)
		if err != nil {
			return err
		}
	}

	for _, d1 := range dirs {
		dir := fs.FileInfoToDirEntry(d1)
		path1 := filepath.Join(path, d1.Name())
		if err := walkDir(repo, path1, dir, walkDirFn); err != nil {
			if err == SkipDir {
				break
			}
			return err
		}
	}
	return nil
}
