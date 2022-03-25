package main

// high-level functional/integration tests

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/gulducat/go-run-programs/program"
	"github.com/stretchr/testify/assert"
)

var setupLock sync.Mutex // hello globe

func TestRunInBackground(t *testing.T) {
	t.Parallel()
	testSetup(t)

	ctx := context.Background()
	p := program.Program{
		Name:    "web0",
		Command: program.Command("test-server start 8080"),
		Check:   program.Command("test-server check 8080"),
	}

	stop, err := program.RunInBackground(ctx, p)
	t.Cleanup(stop)

	if !assert.NoError(t, err) {
		return
	}

	ports := []string{"8080"}
	testHTTP(t, ports)
}

func TestRunFromHCL(t *testing.T) {
	t.Parallel()
	testSetup(t)

	ctx := context.Background()
	// tests config.go ParseHCL too
	stop, err := program.RunFromHCL(ctx, "test_server.hcl")
	t.Cleanup(stop)

	if !assert.NoError(t, err) {
		return
	}

	ports := []string{"8081", "8082", "8083"}
	testHTTP(t, ports)
}

func TestRunInBackground_BadStart(t *testing.T) {
	t.Parallel()
	testSetup(t)

	ctx := context.Background()
	p := program.Program{
		Name:    "BadStart",
		Command: program.Command("test-server death of rats"),
		Check:   program.Command("does-not-matter"),
	}

	_, err := program.RunInBackground(ctx, p)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "of rats")
		assert.Contains(t, err.Error(), "SQUEAK.")
	}
}

func TestRunInBackground_BadCheck(t *testing.T) {
	t.Parallel()
	testSetup(t)

	ctx := context.Background()
	p := program.Program{
		Name:         "BadCheck",
		Command:      program.Command("test-server sleep 60"),
		Check:        program.Command("test-server death of checks"),
		CheckSeconds: 3,
	}

	stop, err := program.RunInBackground(ctx, p)
	t.Cleanup(stop)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "of checks")
		assert.Contains(t, err.Error(), "SQUEAK.")
	}
}

/*** test helpers ***/

func testSetup(t *testing.T) {
	setupLock.Lock()
	defer setupLock.Unlock()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("getwd:", err)
	}
	ensurePathInPATH(t, wd)
	ensureTestServerBin(t, wd)
}

func ensurePathInPATH(t *testing.T, path string) {
	pathSep := ":"
	if runtime.GOOS == "windows" {
		pathSep = ";"
	}
	pathEnv := strings.Trim(os.Getenv("PATH"), pathSep)
	curPaths := strings.Split(pathEnv, pathSep)
	for _, p := range curPaths {
		if p == path {
			return
		}
	}
	paths := append(curPaths, path)
	joined := strings.Join(paths, pathSep)
	t.Log("added "+path+" to PATH:", joined)
	os.Setenv("PATH", joined)
}

func ensureTestServerBin(t *testing.T, path string) {
	file := "test-server"
	if runtime.GOOS == "windows" {
		file += ".exe"
	}
	// if bin already exists, don't spend the time compiling
	if _, err := os.Stat(file); !errors.Is(err, os.ErrNotExist) {
		return
	}

	t.Log("building fakey test-server")
	cmd := exec.Command("go", "build", "-o", file, "./test_server.go")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal("err running command:", err, "; output:", string(out))
	}
	t.Log(string(out))
}

func testHTTP(t *testing.T, ports []string) {
	for _, port := range ports {
		url := "http://127.0.0.1:" + port
		t.Log("test get:", url)
		resp, err := http.Get(url)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	}
}
