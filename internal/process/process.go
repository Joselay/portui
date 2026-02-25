package process

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// Info represents a process listening on a port.
type Info struct {
	Protocol string
	Port     int
	PID      int
	Command  string
	User     string
	State    string
}

// List returns all processes listening on ports.
func List() ([]Info, error) {
	out, err := exec.Command("lsof", "-iTCP", "-iUDP", "-sTCP:LISTEN", "-P", "-n").Output()
	if err != nil {
		// lsof returns exit code 1 when no results found
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, fmt.Errorf("lsof: %w", err)
	}
	return parseLsof(string(out)), nil
}

func parseLsof(output string) []Info {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		return nil
	}

	seen := make(map[string]bool)
	var results []Info

	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		command := fields[0]
		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		user := fields[2]
		protocol := strings.ToUpper(fields[7])
		nameField := fields[8]

		port := parsePort(nameField)
		if port == 0 {
			continue
		}

		state := ""
		if len(fields) > 9 {
			state = strings.Trim(fields[9], "()")
		}

		key := fmt.Sprintf("%s:%d:%d", protocol, port, pid)
		if seen[key] {
			continue
		}
		seen[key] = true

		results = append(results, Info{
			Protocol: protocol,
			Port:     port,
			PID:      pid,
			Command:  command,
			User:     user,
			State:    state,
		})
	}
	return results
}

func parsePort(name string) int {
	// format: *:8080 or 127.0.0.1:8080 or [::1]:8080
	idx := strings.LastIndex(name, ":")
	if idx < 0 {
		return 0
	}
	port, err := strconv.Atoi(name[idx+1:])
	if err != nil {
		return 0
	}
	return port
}

// Kill terminates a process by PID.
func Kill(pid int) error {
	return syscall.Kill(pid, syscall.SIGTERM)
}

// ForceKill sends SIGKILL to a process.
func ForceKill(pid int) error {
	return syscall.Kill(pid, syscall.SIGKILL)
}
