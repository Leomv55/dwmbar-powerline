package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

type CpuLoad struct {
	val map[string]interface{}
}

func (k *CpuLoad) value() map[string]interface{} {
	return k.val
}

func cpu_load() element {
	e := &CpuLoad{}
	go func() {
		for {
			if val, err := e.read(); err == nil {
				e.val = val
			} else {
				log.Printf("could not read cpu load: %v", err)
			}
			time.Sleep(time.Second)
		}
	}()
	return e
}

func (k *CpuLoad) read() (map[string]interface{}, error) {
	data, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, fmt.Errorf("read cpu load from %s - %s", "/proc/loadavg", err)
	}

	parts := strings.Split(strings.TrimSpace(string(data)), " ")
	load, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, fmt.Errorf("parse cpu average load: %s", err)
	}
	var color string
	switch {
	case load >= 10:
		color = "#dc322f"
	case load >= 4:
		color = "#b58900"
	default:
		color = "#6c71c4"
	}
	return map[string]interface{}{
		"color": color,
		"load":  load,
		"icon":  xbm("load")}, nil
}
