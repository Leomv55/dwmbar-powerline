package main

import (
	"log"
	"time"
)

type Date struct {
	val map[string]interface{}
}

func (k *Date) value() map[string]interface{} {
	return k.val
}

func date() element {
	e := &Date{}
	go func() {
		for {
			if val, err := e.read(); val != nil && err == nil {
				e.val = val
			} else {
				log.Printf("could not read date: %v", err)
			}
			time.Sleep(time.Second)
		}
	}()
	return e
}

func (k *Date) read() (map[string]interface{}, error) {
	localTZ, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return nil, err
	}

	local := time.Now().In(localTZ).Format("Mon _2 Jan 15:04")

	return map[string]interface{}{
		"icon": xbm("clock2"),
		"time": local}, nil
}
