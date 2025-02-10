package main

import (
	"errors"
	"math"
	"os/exec"
	logger "github.com/Diamon0/rns-babel/Logger"
	webui "github.com/Diamon0/rns-babel/WebUI"
	"runtime"
	"sync/atomic"

	"github.com/gdamore/tcell/v2"
)

// TODO:
// The addition of the terminal UI to the main package is tentative,
// it should be later decided if it should remain here or moved to its own package

var STYLE_DEFAULT tcell.Style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
var STYLE_BOX tcell.Style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)
var STYLE_BOX_ON tcell.Style = tcell.StyleDefault.Background(tcell.ColorDarkGreen).Foreground(tcell.ColorWhite)
var STYLE_BOX_OFF tcell.Style = tcell.StyleDefault.Background(tcell.ColorDarkRed).Foreground(tcell.ColorWhite)
var STYLE_ON tcell.Style = tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorBlack)
var STYLE_OFF tcell.Style = tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite)

var ServerAddress string = "localhost:3939"

func requestFolderPath() (string, error) {
    var folderPath string

    switch runtime.GOOS {
    case "windows":
    case "darwin":
        cmd := exec.Command("osascript", "-e", `choose folder`)
        output, err := cmd.Output()
        if err != nil {
            return folderPath, err
        }
        folderPath = string(output)

    case "linux":
    default:
        return folderPath, errors.New("Unknown platform")
    }

    return folderPath, nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	// I don't expect these to handle all scenarios,
	// just enough scenarios
	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	default:
		return errors.New("Unknown platform")
	}

	return exec.Command(cmd, args...).Start()
}

type TextAlignment uint8

const (
	UP TextAlignment = iota
	UPRIGHT
	RIGHT
	DOWNRIGHT
	DOWN
	DOWNLEFT
	LEFT
	UPLEFT
	CENTER
)

// Strings (sorta)
const (
	UI_SERVER_POWER    string = "Start Server [S]"
	UI_SERVER_SHUTDOWN string = "Stop Server [K]"
	UI_WEBSITE_OPEN    string = "Open In Browser [O]"
)

// TODO:
// Fix all these draw functions
// This does indeed include recalculating on resize
func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string, textAlignment TextAlignment) {
	textWidth := len([]rune(text))
	boxWidth := x2 - x1 + 1
	boxHeight := y2 - y1 + 1

	col := x1
	row := y1

	switch textAlignment {
	case CENTER:
		col = x1 + (boxWidth-textWidth)/2
		row = y1 + boxHeight/2
	case UP:
		col = x1 + (boxWidth-textWidth)/2
	case DOWN:
		col = x1 + (boxWidth-textWidth)/2
		row = y2
	case LEFT:
		row = y1 + boxHeight/2
	case RIGHT:
		col = x2 - textWidth + 1
		row = y1 + boxHeight/2
	case UPRIGHT:
		col = x2 - textWidth + 1
	case DOWNRIGHT:
		col = x2 - textWidth + 1
		row = y2
	case DOWNLEFT:
		row = y2
	case UPLEFT:
		// Leave as is
	}

	// TODO:
	// I am pretty sure I can make something better than this for handling text.
	// For center:
	// Split into equal parts, maybe based on full size,
	// then split in semi-equal parts based on the next empty space from the 'slice' points
	for _, r := range []rune(text) {
		if col >= x1 && col <= x2 && row >= y1 && row <= y2 {
			s.SetContent(col, row, r, nil, style)
		}
		col++
		if col > x2 {
			col = x1
			row++
		}
		if row > y2 {
			break
		}
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, boxStyle, borderStyle, textStyle tcell.Style, text string, textAlignment TextAlignment) {
	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, boxStyle)
		}
	}

	// Draw Borders
	// Top and Bottom
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, borderStyle)
		s.SetContent(col, y2, tcell.RuneHLine, nil, borderStyle)
	}
	// Left and Right
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, borderStyle)
		s.SetContent(x2, row, tcell.RuneVLine, nil, borderStyle)
	}

	s.SetContent(x1, y1, tcell.RuneULCorner, nil, borderStyle)
	s.SetContent(x2, y1, tcell.RuneURCorner, nil, borderStyle)
	s.SetContent(x1, y2, tcell.RuneLLCorner, nil, borderStyle)
	s.SetContent(x2, y2, tcell.RuneLRCorner, nil, borderStyle)

	drawText(s, x1, y1, x2, y2, textStyle, text, textAlignment)
}

func main() {
	var isServerOn atomic.Bool
	signaler := make(chan int8)
	defer close(signaler)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		logger.DefaultLogger.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		logger.DefaultLogger.Fatalf("%+v", err)
	}
	s.SetStyle(STYLE_DEFAULT)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	// TODO:
	// Do something useful with this
	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	xmax, ymax := s.Size()
	leftBoxXMax := int(math.Floor(float64(xmax)/2)) - 1
	rightBoxXMin := int(math.Ceil(float64(xmax) / 2))
	drawBox(s, 0, 0, leftBoxXMax, ymax-3, STYLE_BOX_OFF, STYLE_BOX, STYLE_ON, UI_SERVER_POWER, CENTER)
	drawBox(s, rightBoxXMin, 0, xmax-1, ymax-3, STYLE_BOX, STYLE_BOX, STYLE_DEFAULT, UI_WEBSITE_OPEN, CENTER)
	drawText(s, 0, ymax-2, xmax, ymax, STYLE_DEFAULT, "Quit [Q or ESC]", UPLEFT)

	// Start main loop
	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Handle event
		switch ev := ev.(type) {
		// If the window was resized
		case *tcell.EventResize:
			//s.Sync()
			s.Clear()
			xmax, ymax = s.Size()
			leftBoxXMax = int(math.Floor(float64(xmax)/2)) - 1
			rightBoxXMin = int(math.Ceil(float64(xmax) / 2))
			if !isServerOn.Load() {
				drawBox(s, 0, 0, leftBoxXMax, ymax-3, STYLE_BOX_OFF, STYLE_BOX, STYLE_ON, UI_SERVER_POWER, CENTER)
			} else {
				drawBox(s, 0, 0, leftBoxXMax, ymax-3, STYLE_BOX_ON, STYLE_BOX, STYLE_OFF, UI_SERVER_SHUTDOWN, CENTER)
			}
			drawBox(s, rightBoxXMin, 0, xmax-1, ymax-3, STYLE_BOX, STYLE_BOX, STYLE_DEFAULT, UI_WEBSITE_OPEN, CENTER)
			drawText(s, 0, ymax-2, xmax, ymax, STYLE_DEFAULT, "Quit [Q or ESC]", UPLEFT)

		// If a key was pressed
		case *tcell.EventKey:
			// Was it for quitting?
			if ev.Key() == tcell.KeyESC || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' || ev.Rune() == 'Q' {
				if isServerOn.Load() {
					signaler <- 1
				}
				return

				// Was it to refresh the display?
			} else if ev.Key() == tcell.KeyCtrlL {
				//s.Sync()
				s.Clear()
				xmax, ymax = s.Size()
				leftBoxXMax = int(math.Floor(float64(xmax)/2)) - 1
				rightBoxXMin = int(math.Ceil(float64(xmax) / 2))
				if !isServerOn.Load() {
					drawBox(s, 0, 0, leftBoxXMax, ymax-3, STYLE_BOX_OFF, STYLE_BOX, STYLE_ON, UI_SERVER_POWER, CENTER)
				} else {
					drawBox(s, 0, 0, leftBoxXMax, ymax-3, STYLE_BOX_ON, STYLE_BOX, STYLE_OFF, UI_SERVER_SHUTDOWN, CENTER)
				}
				drawBox(s, rightBoxXMin, 0, xmax-1, ymax-3, STYLE_BOX, STYLE_BOX, STYLE_DEFAULT, UI_WEBSITE_OPEN, CENTER)
				drawText(s, 0, ymax-2, xmax, ymax, STYLE_DEFAULT, "Quit [Q or ESC]", UPLEFT)

				// Was it to stop the server?
			} else if ev.Rune() == 'k' || ev.Rune() == 'K' {
				if isServerOn.Load() {
					signaler <- 1
					drawBox(s, 0, 0, leftBoxXMax, ymax-3, STYLE_BOX_OFF, STYLE_BOX, STYLE_ON, UI_SERVER_POWER, CENTER)
				}

				// Was it to start the server?
			} else if ev.Rune() == 's' || ev.Rune() == 'S' {
				if !isServerOn.Load() {
					go func() {
						isServerOn.Store(true)
						webui.StartWeb(ServerAddress, signaler)
						isServerOn.Store(false)
					}()
					drawBox(s, 0, 0, leftBoxXMax, ymax-3, STYLE_BOX_ON, STYLE_BOX, STYLE_OFF, UI_SERVER_SHUTDOWN, CENTER)
				}

				// TODO:
				// Consider giving a warning if the server is not on, or maybe preventing it from opening
				//
				// Was it to open the website in the browser?
			} else if ev.Rune() == 'o' || ev.Rune() == 'O' {
				openBrowser("http://" + ServerAddress)
			}

			// If the mouse did something
		case *tcell.EventMouse:
			if ev.Buttons() == tcell.Button1 {
				x, y := ev.Position()

				// Was the left box clicked?
				if x <= leftBoxXMax && y <= ymax-3 {
					// Is the server off?
					if !isServerOn.Load() {
						go func() {
							isServerOn.Store(true)
							webui.StartWeb(ServerAddress, signaler)
							isServerOn.Store(false)
						}()
						drawBox(s, 0, 0, leftBoxXMax, ymax-3, STYLE_BOX_ON, STYLE_BOX, STYLE_OFF, UI_SERVER_SHUTDOWN, CENTER)
					} else {
						// Double-check just in case
						if isServerOn.Load() {
							signaler <- 1
							drawBox(s, 0, 0, leftBoxXMax, ymax-3, STYLE_BOX_OFF, STYLE_BOX, STYLE_ON, UI_SERVER_POWER, CENTER)
						}
					}

					// Was the right box clicked?
				} else if x >= rightBoxXMin && y <= ymax-3 {
					openBrowser("http://" + ServerAddress)
				}
			}
		}
	}
}
