package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"time"
)

type element interface {
	value() map[string]interface{}
}

type statusbar struct {
	elementmap map[string]element

	// configuration properties
	Dzen2               []string        `json:"dzen2"`
	ElementsOrderList   []string        `json:"elements_order_list"`
	DzenPowerlineFormat PowerlineFormat `json:"dzen_powerline_format"`
}

func run(conf string) error {
	var bar statusbar
	file, err := ioutil.ReadFile(conf)
	if err != nil {
		return fmt.Errorf("failed to read config file: %s - %s", conf, err)
	}

	if err := json.Unmarshal(file, &bar); err != nil {
		return fmt.Errorf("failed to unmarshal config file: %s - %s", conf, err)
	}

	bar.elementmap = make(map[string]element, len(bar.ElementsOrderList))
	for _, element := range bar.ElementsOrderList {
		if element == "keyboard" {
			bar.elementmap[element] = keyboard()
		} else if element == "network" {
			bar.elementmap[element] = network()
		} else if element == "cpu_temp" {
			bar.elementmap[element] = cpu_temp()
		} else if element == "cpu_load" {
			bar.elementmap[element] = cpu_load()
		} else if element == "memory_usage" {
			bar.elementmap[element] = memory_usage()
		} else if element == "power" {
			bar.elementmap[element] = power()
		} else if element == "date" {
			bar.elementmap[element] = date()
		} else {
			log.Fatalf("Unknown element type '%s' present", element)
		}
	}
	cmd := exec.Command("dzen2", bar.Dzen2...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %s", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start dzen2 command: %s", err)
	}

	// run the iteration loop
	go func() {
		for {
			if _, e := stdin.Write([]byte(bar.DzenPowerlineFormat.powerlineFormatted(bar.elementmap, bar.ElementsOrderList) + "\n")); e != nil {
				log.Printf("probably the pipe closed: %s", e)
				break
			}
			time.Sleep(time.Second * INTERVAL_SECS)
		}
	}()

	return cmd.Wait()
}
