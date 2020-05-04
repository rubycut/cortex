// +build requires_docker

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/cortexproject/cortex/integration/e2e"
	e2edb "github.com/cortexproject/cortex/integration/e2e/db"
)

var (
	// Expose some utilities from the framework so that we don't have to prefix them
	// with the package name in tests.
	mergeFlags      = e2e.MergeFlags
	newDynamoClient = e2edb.NewDynamoClient
	generateSeries  = e2e.GenerateSeries
)

func getCortexProjectDir() string {
	if dir := os.Getenv("CORTEX_CHECKOUT_DIR"); dir != "" {
		return dir
	}

	return os.Getenv("GOPATH") + "/src/github.com/cortexproject/cortex"
}

func writeFileToSharedDir(s *e2e.Scenario, dst string, content []byte) error {
	dst = filepath.Join(s.SharedDir(), dst)

	// Ensure the entire path of directories exist.
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	return ioutil.WriteFile(
		dst,
		content,
		os.ModePerm)
}

func copyFileToSharedDir(s *e2e.Scenario, src, dst string) error {
	content, err := ioutil.ReadFile(filepath.Join(getCortexProjectDir(), src))
	if err != nil {
		return errors.Wrapf(err, "unable to read local file %s", src)
	}

	return writeFileToSharedDir(s, dst, content)
}

// GetServerTLSFlags generates generic TLS flags for a server
func GetServerTLSFlags() map[string]string {
	return map[string]string{
		"-server.http-tls-cert-path": filepath.Join(e2e.ContainerSharedDir, "certs/server.crt"),
		"-server.http-tls-key-path":  filepath.Join(e2e.ContainerSharedDir, "certs/server.key"),
		"-server.http-tls-ca-path":   filepath.Join(e2e.ContainerSharedDir, "certs/root.crt"),
		"-server.grpc-tls-cert-path": filepath.Join(e2e.ContainerSharedDir, "certs/server.crt"),
		"-server.grpc-tls-key-path":  filepath.Join(e2e.ContainerSharedDir, "certs/server.key"),
		"-server.grpc-tls-ca-path":   filepath.Join(e2e.ContainerSharedDir, "certs/root.crt"),
	}
}

// GetClientTLSFlagsWithPrefix generates generic TLS flags for a client
func GetClientTLSFlagsWithPrefix(prefix string) map[string]string {
	return map[string]string{
		"-" + prefix + ".tls-cert-path": filepath.Join(e2e.ContainerSharedDir, "certs/client.crt"),
		"-" + prefix + ".tls-key-path":  filepath.Join(e2e.ContainerSharedDir, "certs/client.key"),
		"-" + prefix + ".tls-ca-path":   filepath.Join(e2e.ContainerSharedDir, "certs/root.crt"),
	}
}
