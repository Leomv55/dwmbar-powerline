package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var match_mem = regexp.MustCompile(`Mem:\s+(.+)`)

type Memory struct {
	val map[string]interface{}
}

func (k *Memory) value() map[string]interface{} {
	return k.val
}

func memory_usage() element {
	e := &Memory{}
	go func() {
		for {
			if val, err := e.read(); err == nil {
				e.val = val
			} else {
				log.Printf("could not read memory usage: %v", err)
			}
			time.Sleep(time.Second)
		}
	}()
	return e
}

func (k *Memory) read() (map[string]interface{}, error) {
	data, err := exec.Command("free", "-m").Output()
	if err != nil {
		return nil, fmt.Errorf("'free -m' command: %s", err)
	}

	m := match_mem.FindStringSubmatch(string(data))
	if len(m) != 2 {
		return nil, fmt.Errorf("number of matches was not expected for mem submatch")
	}

	parts := strings.Fields(m[1])
	total, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("read total memory: %s", err)
	}
	avail, err := strconv.Atoi(parts[5])
	if err != nil {
		return nil, fmt.Errorf("read available memory: %s", err)
	}
	used := total - avail
	perc := 100 * used / total

	var color string
	switch {
	case perc >= 90:
		color = "#dc322f"
	case perc >= 70:
		color = "#b58900"
	default:
		color = "#859900"
	}

	return map[string]interface{}{
		"color": color,
		"perc":  perc,
		"icon":  xbm("mem")}, nil
}
