package file

import (
	"io/ioutil"
	nurl "net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/neighborly/ddsl/drivers/source"
)

func init() {
	source.Register("file", &File{})
}

type File struct {
	url        string
	path       string
}

func Register() {
	// do nothing, but call to force compiler to accept import without use
}

func (f *File) Open(url string) (source.Driver, error) {
	u, err := nurl.Parse(url)
	if err != nil {
		return nil, err
	}

	// concat host and path to restore full path
	// host might be `.`
	p := u.Opaque
	if len(p) == 0 {
		p = u.Host + u.Path
	}

	if len(u.Fragment) > 0 {
		p += "#" + u.Fragment
	}

	if len(p) == 0 {
		// default to current directory if no path
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		p = wd

	} else if p[0:1] == "." || p[0:1] != "/" {
		// make path absolute if relative
		abs, err := filepath.Abs(p)
		if err != nil {
			return nil, err
		}
		p = abs
	}

	nf := &File{
		url:        url,
		path:       p,
	}

	return nf, nil
}

func (f *File) Close() error {
	// nothing do to here
	return nil
}

func (f *File) readDirectory(relativeDir string, fileNamePattern string, recursive bool) (df *source.DirectoryReader, err error) {
	dirPath := path.Join(f.path, relativeDir)
	var re *regexp.Regexp
	if len(fileNamePattern) > 0 {
		re = regexp.MustCompile(fileNamePattern)
	}

	dr := &source.DirectoryReader{
		DirectoryPath: dirPath,
		FileReaders: []*source.FileReader{},
		SubDirectories: []*source.DirectoryReader{},
	}

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return dr, nil
	}

	items, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		itemPath := path.Join(dirPath, item.Name())
		if item.IsDir() {
			var subdr *source.DirectoryReader
			if recursive {
				if subdr, err = f.readDirectory(itemPath, fileNamePattern, recursive); err != nil {
					return nil, err
				}
			} else {
				subdr = &source.DirectoryReader{
					DirectoryPath: path.Join(dirPath, item.Name()),
				}
			}
			dr.SubDirectories = append(dr.SubDirectories, subdr)
		} else {
			match := re == nil || re.MatchString(item.Name())
			if !match {
				continue
			}

			reader, err := os.Open(itemPath)
			if err != nil {
				return nil, err
			}
			fr := &source.FileReader{
				FilePath: itemPath,
				Reader: reader,
			}
			dr.FileReaders = append(dr.FileReaders, fr)
		}
	}

	return dr, nil
}

func (f *File) ReadFiles(relativeDir string, fileNamePattern string) (files []*source.FileReader, err error) {
	dr, err := f.readDirectory(relativeDir, fileNamePattern, false)
	if err != nil {
		return nil, err
	}
	return dr.FileReaders, nil
}

func (f *File) ReadDirectories(relativeDir string, dirNamePattern string) (files []*source.DirectoryReader, err error) {
	dr, err := f.readDirectory(relativeDir, dirNamePattern, false)
	if err != nil {
		return nil, err
	}
	return dr.SubDirectories, nil
}

func (f *File) ReadTree(relativeDir string, fileNamePattern string) (tree *source.DirectoryReader, err error) {
	return f.readDirectory(relativeDir, fileNamePattern, true)
}

