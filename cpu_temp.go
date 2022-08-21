package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var match_temp = regexp.MustCompile(`Core\s+\d+:\s+\+([\d]+)`)

type CpuTemp struct {
	val    map[string]interface{}
	inputs []string
}

func (k *CpuTemp) value() map[string]interface{} {
	return k.val
}

func cpu_temp() element {
	e := &CpuTemp{}
	// collect temp inputs, this may vary based on different kernels
	for i := 1; i < 128; i++ { // might have 128 cores?
		p := fmt.Sprintf("/sys/class/hwmon/hwmon0/temp%d_input", i)
		if file_exists(p) {
			e.inputs = append(e.inputs, p)
		} else {
			break
		}
	}
	if len(e.inputs) == 0 {
		// look elsewhere
		p := "/sys/class/thermal/thermal_zone0/temp"
		if file_exists(p) {
			e.inputs = append(e.inputs, p)
		}
	}
	if len(e.inputs) == 0 {
		log.Println("could not determine CPU temperature inputs")
		return e
	}
	go func() {
		for {
			if val, err := e.read(); err == nil {
				e.val = val
			} else {
				log.Printf("could not read cpu temp: %v", err)
			}
			time.Sleep(time.Second * 2)
		}
	}()
	return e
}

func (k *CpuTemp) read() (map[string]interface{}, error) {
	var avg int
	for _, in := range k.inputs {
		data, err := ioutil.ReadFile(in)
		if err != nil {
			return nil, fmt.Errorf("failed to read temperature input: %s - %v", in, err)
		}

		s := strings.TrimSpace(string(data))
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("failed to read temperature input: %s - %v", in, err)
		}
		avg += i
	}

	temp := avg / len(k.inputs) / 1000
	var color string
	switch {
	case temp >= 80:
		color = "#dc322f"
	case temp >= 60:
		color = "#b58900"
	default:
		color = "#859900"
	}
	return map[string]interface{}{
		"color": color,
		"temp":  temp}, nil
}
