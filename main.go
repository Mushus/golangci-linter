package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/golangci/golangci-lint/pkg/result"
)

const (
	envGitHubWorkspace = "GITHUB_WORKSPACE"
	envConfig          = "INPUT_CONFIG"
	envBasePath        = "INPUT_BASEPATH"
)

type report struct {
	Issues []result.Issue `json:"Issues"`
	Report struct {
		Error string `json:"Error"`
	} `json:"Report"`
}

var typecheckingErrRegexp = regexp.MustCompile(`^typechecking\serror:\s(.+\.go):(\d+):(\d+):\s(.+)$`)

func (r report) parseError() []annotation {
	anns := make([]annotation, 0)
	m := typecheckingErrRegexp.FindStringSubmatch(r.Report.Error)
	if len(m) == 0 {
		return anns
	}
	line, err := strconv.Atoi(m[2])
	if err != nil {
		// not fail
		return anns
	}
	col, err := strconv.Atoi(m[3])
	if err != nil {
		// not fail
		return anns
	}
	anns = append(anns, annotation{
		file: m[1],
		line: line,
		col:  col,
		text: m[4],
	})
	return anns
}

type annotation struct {
	file string
	line int
	col  int
	text string
}

func (a annotation) Output() string {
	return fmt.Sprintf("::error file=%s,line=%d,col=%d::%s", a.file, a.line, a.col, a.text)
}

type config struct {
	workspace string
	config    string
	basePath  string
}

func loadConfig() config {
	return config{
		workspace: os.Getenv(envGitHubWorkspace),
		config:    os.Getenv(envConfig),
		basePath:  os.Getenv(envBasePath),
	}
}

func decodeJSON(r io.Reader) (*report, error) {
	var rep report
	if err := json.NewDecoder(r).Decode(&rep); err != nil {
		return nil, fmt.Errorf("cannot parse result: %w", err)
	}
	return &rep, nil
}

func execGolangCILint(cfg config) (int, []annotation, error) {
	args := []string{"run"}
	if cfg.config != "" {
		args = append(args, "--config", cfg.config)
	}
	args = append(args, "--out-format", "json")

	baseDir := filepath.Join(cfg.workspace, cfg.basePath)
	fmt.Printf("%v\n", args)
	fmt.Printf("%v %v %v\n", baseDir, cfg.workspace, cfg.basePath)
	cmd0 := exec.Command("ls", "-al")
	cmd0.Dir = baseDir
	out, err := cmd0.Output()
	fmt.Printf("%v\n", err)
	fmt.Printf("%v\n", out)
	cmd := exec.Command("golangci-lint", args...)
	cmd.Dir = baseDir
	r, err := cmd.StdoutPipe()
	if err != nil {
		return -1, nil, fmt.Errorf("cannot get stdout: %w", err)
	}
	defer r.Close()
	if err := cmd.Start(); err != nil {
		return -1, nil, fmt.Errorf("cannot execute lint: %w", err)
	}

	b, err := ioutil.ReadAll(r)
	fmt.Printf("%v\n", err)
	fmt.Printf("%s\n", b)
	r2 := bytes.NewReader(b)

	var rep report
	if err := json.NewDecoder(r2).Decode(&rep); err != nil {
		return -1, nil, fmt.Errorf("cannot parse result: %w", err)
	}

	// NOTE: Not need error checking.
	// When `cmd.Wait()` failed, `cmd` has already completed.
	//nolint:errcheck
	cmd.Wait()
	exitCode := cmd.ProcessState.ExitCode()

	if errAnns := rep.parseError(); len(errAnns) > 0 {
		return exitCode, errAnns, nil
	}

	anns := createAnotations(cfg, rep.Issues)

	return exitCode, anns, nil
}

func createAnotations(cfg config, issues []result.Issue) []annotation {
	ann := make([]annotation, len(issues))
	for i := range issues {
		pos := issues[i].Pos
		file := filepath.Join(cfg.basePath, pos.Filename)
		ann[i] = annotation{
			file: file,
			line: pos.Line,
			col:  pos.Column,
			text: fmt.Sprintf("[%s] %s", issues[i].FromLinter, issues[i].Text),
		}
	}
	return ann
}

func reportFailures(cfg config, anns []annotation) {
	for _, ann := range anns {
		fmt.Println(ann.Output())
	}
}

func main() {
	cfg := loadConfig()

	exitStatus, issues, err := execGolangCILint(cfg)
	if err != nil {
		log.Fatalf("failed to execute golangci-lint: %v", err)
	}

	if len(issues) > 0 {
		reportFailures(cfg, issues)
	}

	os.Exit(exitStatus)
}
