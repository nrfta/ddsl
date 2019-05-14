// Package source provides the Source interface.
// All source drivers must implement this interface, register themselves,
// optionally provide a `WithInstance` function and pass the tests
// in package source/testing.
package source

import (
	"fmt"
	"io"
	nurl "net/url"
	"sync"
)

var driversMu sync.RWMutex
var drivers = make(map[string]Driver)

// Driver is the interface every source driver must implement.
//
// How to implement a source driver?
//   1. Implement this interface.
//   2. Optionally, add a function named `WithInstance`.
//      This function should accept an existing source instance and a Config{} struct
//      and return a driver instance.
//   3. Add a test that calls source/testing.go:Test()
//   4. Add own tests for Open(), WithInstance() (when provided) and Close().
//      All other functions are tested by tests in source/testing.
//      Saves you some time and makes sure all source drivers behave the same way.
//   5. Call Register in init().
//
// Guidelines:
//   * All configuration input must come from the URL string in func Open()
//     or the Config{} struct in WithInstance. Don't os.Getenv().
//   * Drivers are supposed to be read only.
//   * Ideally don't load any contents (into memory) in Open or WithInstance.
type Driver interface {
	// Open returns a a new driver instance configured with parameters
	// coming from the URL string.
	Open(url string) (Driver, error)

	// Close closes the underlying source instance managed by the driver.
	Close() error

	// ReadFiles returns `FileReader` slice for the file at the given relative directory
	// that match the given pattern. Returns `nil` if relative path does not exist.
	ReadFiles(relativeDir string, fileNamePattern string) (files []*FileReader, err error)

	// ReadDirectories returns `DirectoryReader` slice for the directories at the given relative directory
	// that match the given pattern. Returns `nil` if relative path does not exist.
	ReadDirectories(relativeDir string, dirNamePattern string) (files []*DirectoryReader, err error)

	// ReadTree returns a `DirectoryReader` for directory at the given relative path
	// with the `SubDirectories` member recursively populated. Returns `nil` if the path
	// does not exist.
	ReadTree(relativeDir string, fileNamePattern string) (tree *DirectoryReader, err error)
}

type FileReader struct {
	FilePath string
	Reader io.ReadCloser
}

type DirectoryReader struct {
	DirectoryPath string
	FileReaders []*FileReader
	SubDirectories []*DirectoryReader
}

// Open returns a new driver instance.
func Open(url string) (Driver, error) {
	u, err := nurl.Parse(url)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		return nil, fmt.Errorf("source driver: invalid URL scheme")
	}

	driversMu.RLock()
	d, ok := drivers[u.Scheme]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("source driver: unknown driver %v (forgotten import?)", u.Scheme)
	}

	return d.Open(url)
}

// Register globally registers a driver.
func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("Register called twice for driver " + name)
	}
	drivers[name] = driver
}

// List lists the registered drivers
func List() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	names := make([]string, 0, len(drivers))
	for n := range drivers {
		names = append(names, n)
	}
	return names
}