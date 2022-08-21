package main

import (
	"fmt"
	"log"
)

type PowerlineNetwokFormat struct {
	Disconnected   string `json:"disconnected"`
	WifiFull       string `json:"wifi_full"`
	WifiMid        string `json:"wifi_mid"`
	WifiLow        string `json:"wifi_low"`
	WiredConnected string `json:"wired_connected"`
}

func (networkStatus *PowerlineNetwokFormat) getPowerlineStringForStatus(
	status string) string {
	var pwNetworkFormat string
	switch status {
	case "wired_connected":
		pwNetworkFormat = networkStatus.WiredConnected
	case "wifi_full":
		pwNetworkFormat = networkStatus.WifiFull
	case "wifi_mid":
		pwNetworkFormat = networkStatus.WifiMid
	case "wifi_low":
		pwNetworkFormat = networkStatus.WifiLow
	default:
		pwNetworkFormat = networkStatus.Disconnected
	}

	return pwNetworkFormat
}

type PowerlineFormat struct {
	Keyboard string                `json:"keyboard"`
	CpuTemp  string                `json:"cpu_temp"`
	CpuLoad  string                `json:"cpu_load"`
	Memory   string                `json:"memory"`
	Power    string                `json:"power"`
	Date     string                `json:"date"`
	Network  PowerlineNetwokFormat `json:"network"`
}

func (powerlineconfig *PowerlineFormat) powerlineFormatted(elementmap map[string]element, elementOrder []string) string {

	keyboardFormat := func() string {
		var powerlineFormated string
		layout := elementmap["keyboard"].value()["layout"]
		if layout == nil {
			powerlineFormated = ""
		} else {
			powerlineFormated = fmt.Sprintf(powerlineconfig.Keyboard, layout.(string))
		}
		return powerlineFormated
	}

	networkFormat := func() string {
		el := elementmap["network"].value()
		networkStatus := el["status"].(string)
		powerlineNetworkFormat := powerlineconfig.Network.getPowerlineStringForStatus(
			networkStatus)
		powerlineFormated := fmt.Sprintf(powerlineNetworkFormat,
			el["icon"].(string), el["traffic"].(string))
		return powerlineFormated
	}

	cpuTempFormat := func() string {
		el := elementmap["cpu_temp"].value()
		var powerlineFormated string
		if el == nil {
			powerlineFormated = ""
		} else {
			powerlineFormated = fmt.Sprintf(powerlineconfig.CpuTemp,
				el["color"].(string), el["temp"].(int))
		}
		return powerlineFormated
	}

	cpuLoadFormat := func() string {
		el := elementmap["cpu_load"].value()
		var powerlineFormated string

		if el == nil {
			powerlineFormated = ""
		} else {
			powerlineFormated = fmt.Sprintf(powerlineconfig.CpuLoad,
				el["color"].(string), el["load"].(float64), el["icon"].(string))
		}
		return powerlineFormated
	}
	powerFormat := func() string {
		var powerlineFormated string
		el := elementmap["power"].value()
		percentage := el["perc"]
		color := el["color"]
		icon := el["icon"]

		if percentage == nil || color == nil || icon == nil {
			powerlineFormated = ""
		} else {
			powerlineFormated = fmt.Sprintf(powerlineconfig.Power,
				color.(string), percentage.(int), icon.(string))

		}
		return powerlineFormated
	}

	memoryUsageFormat := func() string {
		el := elementmap["memory_usage"].value()
		var powerlineFormated string

		if el == nil {
			powerlineFormated = ""
		} else {
			powerlineFormated = fmt.Sprintf(powerlineconfig.Memory,
				el["color"].(string), el["perc"].(int), el["icon"].(string))
		}
		return powerlineFormated
	}

	timeFormat := func() string {
		el := elementmap["date"].value()
		var powerlineFormated string
		if el == nil {
			powerlineFormated = ""
		} else {
			powerlineFormated = fmt.Sprintf(powerlineconfig.Date,
				el["icon"].(string), el["time"].(string))
		}
		return powerlineFormated
	}

	var powelineFormattedStrings string

	for _, k := range elementOrder {

		if k == "keyboard" {
			powelineFormattedStrings += keyboardFormat()
		} else if k == "network" {
			powelineFormattedStrings += networkFormat()
		} else if k == "cpu_temp" {
			powelineFormattedStrings += cpuTempFormat()
		} else if k == "power" {
			powelineFormattedStrings += powerFormat()
		} else if k == "cpu_load" {
			powelineFormattedStrings += cpuLoadFormat()
		} else if k == "memory_usage" {
			powelineFormattedStrings += memoryUsageFormat()
		} else if k == "date" {
			powelineFormattedStrings += timeFormat()
		} else {
			log.Fatal()
		}
		powelineFormattedStrings += " "
	}
	return powelineFormattedStrings
}
