package onprem

import (
	"ivory/plugins/platform"
	"path"
	"strconv"
	"strings"
)

func (a *Adapter) Processes(connection platform.Connection) ([]platform.Process, error) {
	result, err := a.execute(connection, ProcessesCommand)
	if err != nil {
		return nil, err
	}
	return a.parseProcesses(result)
}

// parseProcesses is lenient: a single malformed row (e.g. a command name with
// an unexpected shape) is skipped rather than failing the whole listing.
func (a *Adapter) parseProcesses(lines []string) ([]platform.Process, error) {
	processes := make([]platform.Process, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		fields := strings.Fields(trimmed)
		if len(fields) < 8 {
			continue
		}

		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		cpuPercent, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			continue
		}
		memPercent, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			continue
		}
		rssKb, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			continue
		}
		threads, err := strconv.Atoi(fields[5])
		if err != nil {
			continue
		}

		processes = append(processes, platform.Process{
			Pid:         pid,
			Program:     path.Base(fields[6]),
			Command:     strings.Join(fields[7:], " "),
			Threads:     threads,
			User:        fields[1],
			MemoryBytes: rssKb * 1024,
			MemPercent:  memPercent,
			CpuPercent:  cpuPercent,
		})
	}

	if len(processes) == 0 && len(lines) > 0 {
		return nil, platform.ErrInvalidProcesses
	}

	return processes, nil
}
