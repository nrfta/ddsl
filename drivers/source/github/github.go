package github

import (
	"context"
	"fmt"
	"io/ioutil"
	nurl "net/url"
	"path"
	"regexp"
	"strings"
)

import (
	"github.com/google/go-github/github"
	"github.com/neighborly/ddsl/drivers/source"
)

func init() {
	source.Register("github", &Github{})
}

var (
	ErrNoUserInfo          = fmt.Errorf("no username:token provided")
	ErrNoAccessToken       = fmt.Errorf("no access token")
	ErrInvalidRepo         = fmt.Errorf("invalid repo")
	ErrInvalidGithubClient = fmt.Errorf("expected *github.Client")
	ErrNoDir               = fmt.Errorf("no directory")
)

type Github struct {
	client *github.Client
	url    string

	pathOwner  string
	pathRepo   string
	path       string
	options    *github.RepositoryContentGetOptions
}


type Config struct {
}

func (g *Github) Open(url string) (source.Driver, error) {
	u, err := nurl.Parse(url)
	if err != nil {
		return nil, err
	}

	if u.User == nil {
		return nil, ErrNoUserInfo
	}

	password, ok := u.User.Password()
	if !ok {
		return nil, ErrNoUserInfo
	}

	tr := &github.BasicAuthTransport{
		Username: u.User.Username(),
		Password: password,
	}

	gn := &Github{
		client:     github.NewClient(tr.Client()),
		url:        url,
		options:    &github.RepositoryContentGetOptions{Ref: u.Fragment},
	}

	// set owner, repo and path in repo
	gn.pathOwner = u.Host
	pe := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(pe) < 1 {
		return nil, ErrInvalidRepo
	}
	gn.pathRepo = pe[0]
	if len(pe) > 1 {
		gn.path = strings.Join(pe[1:], "/")
	}

	return gn, nil
}

func WithInstance(client *github.Client, config *Config) (source.Driver, error) {
	gn := &Github{
		client:     client,
	}

	return gn, nil
}

func (g *Github) readDirectory(relativeDir string, fileNamePattern string, recursive bool) (*source.DirectoryReader, error) {
	dirPath := path.Join(g.path, relativeDir)
	fileContent, dirContents, _, err := g.client.Repositories.GetContents(context.Background(), g.pathOwner, g.pathRepo, dirPath, g.options)
	if err != nil {
		return nil, err
	}
	if fileContent != nil {
		return nil, ErrNoDir
	}

	var re *regexp.Regexp
	if len(fileNamePattern) > 0 {
		re = regexp.MustCompile(fileNamePattern)
	}

	dr := &source.DirectoryReader{
		DirectoryPath: dirPath,
		FileReaders: []*source.FileReader{},
		SubDirectories: []*source.DirectoryReader{},
	}

	for _, item := range dirContents {
		name := item.GetName()
		switch item.GetType() {
		case "dir":
			var subdr *source.DirectoryReader
			subDirPath := path.Join(relativeDir, name)
			if recursive {
				subdr, err = g.readDirectory(subDirPath, fileNamePattern, recursive)
				if err != nil {
					return nil, err
				}
			} else {
				subdr = &source.DirectoryReader{DirectoryPath: subDirPath}
			}
			dr.SubDirectories = append(dr.SubDirectories, subdr)

		case "file":
			match := re == nil || re.MatchString(name)
			if !match {
				continue // ignore files that we can't parse
			}

			r, err := item.GetContent()
			if err != nil {
				return nil, err
			}

			fr := &source.FileReader{
				Reader:   ioutil.NopCloser(strings.NewReader(r)),
				FilePath: path.Join(relativeDir, name),
			}

			dr.FileReaders = append(dr.FileReaders, fr)
		}
	}

	return dr, nil
}

func (g *Github) Close() {
	return nil
}

func (g *Github) ReadFiles(relativeDir string, fileNamePattern string) (files []*source.FileReader, err error) {
	dr, err := g.readDirectory(relativeDir, fileNamePattern, false)
	if err != nil {
		return nil, err
	}

	return dr.FileReaders, nil
}


func (g *Github) ReadTree(relativePath string, fileNamePattern string) (t *source.DirectoryReader, err error) {
	return g.readDirectory(relativePath, fileNamePattern, true)
}

