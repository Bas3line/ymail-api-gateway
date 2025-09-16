package utils

import (
	"log"
	"runtime"
)

func SafeGoroutine(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Goroutine panic recovered: %v", r)
				for i := 0; i < 10; i++ {
					if pc, file, line, ok := runtime.Caller(i); ok {
						log.Printf("  %s:%d %s", file, line, runtime.FuncForPC(pc).Name())
					}
				}
			}
		}()
		fn()
	}()
}