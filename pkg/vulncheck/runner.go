package vulncheck

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/vuln/vulncheck"
)

const (
	command    = "govulncheck"
	flag       = "-json"
	envPackage = "PACKAGE"
)

type Scanner interface {
	Scan() (*vulncheck.Result, error)
}

type CmdScanner struct {
}

func NewScanner() Scanner {
	return &CmdScanner{}
}

func (r *CmdScanner) Scan() (*vulncheck.Result, error) {
	pkg := os.Getenv(envPackage)
	workDir, _ := os.Getwd()

	fmt.Printf("Running govulncheck for package %s in dir %s\n", pkg, workDir)
	cmd := exec.Command(command, flag, pkg)
	cmd.Dir = workDir

	out, cmdErr := cmd.Output()
	if err, ok := cmdErr.(*exec.ExitError); ok {
		if err.ExitCode() > 0 {
			println("Scan found vulnerabilities in codebase")
		}

		if len(err.Stderr) > 0 {
			fmt.Printf("Stderr: %s\n", string(err.Stderr))
			fmt.Printf("Error: %v\n", err)
		}

	} else if cmdErr != nil {
		return nil, cmdErr
	}

	var result vulncheck.Result
	err := json.Unmarshal(out, &result)
	if err != nil {
		return nil, errors.New("scan failed to produce proper report")
	}

	fmt.Println("Successfully parsed report")
	return &result, nil
}
