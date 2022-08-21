package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type pw_sources struct {
	AC        string
	val       map[string]interface{}
	batteries []string
}

func (s *pw_sources) value() map[string]interface{} {
	return s.val
}

func power() element {
	e := &pw_sources{}
	if err := e.prepare(); err != nil {
		log.Printf("failed to prepare battery: %v\n", err)
		return e
	}

	go func() {
		for {
			e.val = e.read()
			time.Sleep(time.Second * 3)
		}
	}()
	return e
}

func (s *pw_sources) prepare() error {
	devs, err := ioutil.ReadDir("/sys/class/power_supply")
	if err != nil {
		return err
	}

	for _, dev := range devs {
		d := filepath.Base(dev.Name())
		p := filepath.Join("/sys/class/power_supply", d)
		// filter out non devices
		if !file_exists(filepath.Join(p, "device")) {
			continue // not a physical device
		}

		// maybe battery
		if strings.Index(d, "BAT") != -1 {
			cap := filepath.Join(p, "capacity")
			if !file_exists(cap) {
				return fmt.Errorf("could not locate battery capacity stats at: %s", cap)
			}
			s.batteries = append(s.batteries, cap)
		}

		if d == "AC" {
			s.AC = filepath.Join(p, "online")
			if !file_exists(s.AC) {
				return fmt.Errorf("could not locate AC online stat at: %s", s.AC)
			}
		}
	}

	return nil
}

func (s *pw_sources) onAC() bool {
	if len(s.batteries) == 0 {
		return true
	}

	if len(s.AC) == 0 {
		return false // should not be the case
	}

	dat, err := ioutil.ReadFile(s.AC)
	if err != nil {
		return false
	}

	if strings.TrimSpace(string(dat)) == "0" {
		return false
	}

	return true
}

func (s *pw_sources) battery() int {
	if len(s.batteries) == 0 {
		return 0 // no baterries
	}

	var all int
	for _, b := range s.batteries {
		dat, err := ioutil.ReadFile(b)
		if err != nil {
			continue
		}
		i, _ := strconv.Atoi(strings.TrimSpace(string(dat)))
		all += i
	}

	return all / len(s.batteries)
}

func (s *pw_sources) read() map[string]interface{} {
	if s.onAC() {
		return map[string]interface{}{
			"icon":  xbm("power-ac"),
			"perc":  nil,
			"color": nil,
		}
	}

	perc := s.battery()
	var color, icon string
	switch {
	case perc <= 20:
		icon = xbm("bat-low")
		color = "#dc322f"
	case perc <= 50:
		icon = xbm("bat-mid")
		color = "#b58900"
	default:
		icon = xbm("bat-full")
		color = "#859900"
	}

	return map[string]interface{}{
		"icon":  icon,
		"perc":  perc,
		"color": color}
}
