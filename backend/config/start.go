package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type AnalyzerStart struct {
	cmd *exec.Cmd
	Via string
}

func existingFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func uniqueDirs(paths ...string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		if strings.TrimSpace(p) == "" {
			continue
		}
		clean := filepath.Clean(p)
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		out = append(out, clean)
	}
	return out
}

func StartDisplayAnalyzer(exeDir string) (*AnalyzerStart, error) {
	cwd, _ := os.Getwd()
	searchDirs := uniqueDirs(
		exeDir,
		cwd,
		filepath.Dir(cwd),
		filepath.Join(exeDir, ".."),
	)

	binaryNames := []string{"displayAnalyzer.exe", "displayAnalyzer"}
	for _, dir := range searchDirs {
		for _, name := range binaryNames {
			direct := filepath.Join(dir, name)
			nested := filepath.Join(dir, "displayAnalyzer", "dist", name)
			for _, candidate := range []string{direct, nested} {
				if !existingFile(candidate) {
					continue
				}
				cmd := exec.Command(candidate)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Start(); err != nil {
					log.Printf("Analyzer candidate %q failed to start: %v", candidate, err)
					continue
				}
				go func() {
					if err := cmd.Wait(); err != nil {
						log.Printf("Display analyzer exited with error: %v", err)
						return
					}
					log.Printf("Display analyzer stopped.")
				}()
				return &AnalyzerStart{
					cmd: cmd,
					Via: candidate,
				}, nil
			}
		}
	}

	pyScriptCandidates := make([]string, 0, len(searchDirs))
	for _, dir := range searchDirs {
		pyScriptCandidates = append(pyScriptCandidates, filepath.Join(dir, "displayAnalyzer", "main.py"))
	}

	pythonCommands := [][]string{
		{"python3"},
		{"python"},
	}
	if runtime.GOOS == "windows" {
		pythonCommands = append([][]string{{"py", "-3"}}, pythonCommands...)
	}

	for _, script := range pyScriptCandidates {
		if !existingFile(script) {
			continue
		}
		for _, baseCmd := range pythonCommands {
			if _, err := exec.LookPath(baseCmd[0]); err != nil {
				continue
			}
			args := append(baseCmd[1:], script)
			cmd := exec.Command(baseCmd[0], args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Dir = filepath.Dir(script)
			if err := cmd.Start(); err != nil {
				return nil, fmt.Errorf("failed starting python analyzer via %q %q: %w", baseCmd[0], script, err)
			}
			go func() {
				if err := cmd.Wait(); err != nil {
					log.Printf("Display analyzer exited with error: %v", err)
					return
				}
				log.Printf("Display analyzer stopped.")
			}()
			return &AnalyzerStart{
				cmd: cmd,
				Via: fmt.Sprintf("%s %s", baseCmd[0], script),
			}, nil
		}
	}

	return nil, errors.New("could not find analyzer binary or python script")
}
