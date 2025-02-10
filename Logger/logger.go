package logger

import (
	"log"
	"os"
	"sync"
)

// So, what happens if someone opens the program twice at the same time? (Diamon).
// UHGG, don't make me thing about it, let's imagine it won't happen (Also Diamon).
type LogMutex struct {
	Logger *log.Logger
	M      sync.Mutex
}

func (lm *LogMutex) Println(v ...any) {
	lm.M.Lock()
	lm.Logger.Println(v...)
	lm.M.Unlock()
}

func (lm *LogMutex) Fatalf(format string, v ...any) {
	lm.M.Lock()
	lm.Logger.Fatalf(format, v...)
	lm.M.Unlock()
}

var DefaultLogger LogMutex

func init() {
	// TODO:
	// Check what caused the error, ignore if the file just doesn't exist
	if err := os.Remove("logs.txt"); err != nil {
	}
	logfile, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to initialize program logger\n%+v", err)
	}
	DefaultLogger = LogMutex{
		// TODO: (Maybe?)
		// Later change it so it depends on what module the message came from
		Logger: log.New(logfile, "[Logger]", 0),
		M:      sync.Mutex{},
	}

}
