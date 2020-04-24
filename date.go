package main

import (
	"fmt"
	"log"
	"time"
)

type Date struct {
	val string
}

func (k *Date) value() string {
	return k.val
}

func date() element {
	e := &Date{}
	go func() {
		for {
			if val, err := e.read(); err == nil {
				e.val = val
			} else {
				log.Printf("could not read date: %v", err)
			}
			time.Sleep(time.Second)
		}
	}()
	return e
}

func (k *Date) read() (string, error) {
	localTZ, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return "", err
	}
	

	local := time.Now().In(localTZ).Format("Mon _2 Jan 15:04")

	return fmt.Sprintf("^fg(white)| ^i(%s) ^fg(white)%s ",  xbm("clock2"), local), nil
}
