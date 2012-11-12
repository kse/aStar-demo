package main

import (
	"github.com/banthar/Go-SDL/sdl"
	"fmt"
)

const (
	OPEN = iota
	WALL
)

func fillBoxWall(w [][]int, screen *sdl.Surface, r *sdl.Rect) {
	w[(r.X/20)][(r.Y/20)] = WALL
	screen.FillRect(r, 0xFF4040)
	screen.UpdateRect(
		int32(r.X),
		int32(r.Y),
		uint32(r.W),
		uint32(r.H))
	}

func fillBoxOpen(w [][]int, screen *sdl.Surface, r *sdl.Rect) {
	w[(r.X/20)][(r.Y/20)] = OPEN
	screen.FillRect(r, 0x33CCCC)
	screen.UpdateRect(
		int32(r.X),
		int32(r.Y),
		uint32(r.W),
		uint32(r.H));
}

func drawLine(w [][]int, screen *sdl.Surface, e *sdl.MouseMotionEvent) {
	var x1, y1 int = int(e.X) - int(e.Xrel), int(e.Y) - int(e.Yrel)
	var x0, y0 int = int(e.X), int(e.Y)

	dx := abs(x1 - x0)
	dy := abs(y1 - y0)

	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}

	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}

	var err, e2 int
	if dx > dy {
		err = dx/2
	} else {
		err = -dy/2
	}

	for true {
		r := &sdl.Rect{
				int16(x0 - (x0%20) + 1),
				int16(y0 - (y0%20) + 1),
				19,
				19}

		fillBoxWall(w, screen, r)

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 = err

		if e2 > -dx {
			err -= dy
			x0 += sx
		}

		if e2 < dy {
			err += dx
			y0 += sy
		}
	}
	/*
	void line(int x0, int y0, int x1, int y1) {

		int dx = abs(x1-x0), sx = x0<x1 ? 1 : -1;
		int dy = abs(y1-y0), sy = y0<y1 ? 1 : -1; 
		int err = (dx>dy ? dx : -dy)/2, e2;

		for(;;){
			setPixel(x0,y0);
			if (x0==x1 && y0==y1) break;
			e2 = err;
			if (e2 >-dx) { err -= dy; x0 += sx; }
			if (e2 < dy) { err += dx; y0 += sy; }
		}
	}
	*/
}

func abs(v int) int{
	if v < 0 {
		return v * -1
	}

	return v
}

func main() {
	// Contains our world, which is simply an array of types
	var world [][]int

	if sdl.Init(sdl.INIT_VIDEO) != 0 {
		panic(sdl.GetError())
	}

	v_info := sdl.GetVideoInfo()

	var screen = sdl.SetVideoMode(
		int(v_info.Current_w),
		int(v_info.Current_h),
		32,
		sdl.HWSURFACE | sdl.DOUBLEBUF | sdl.FULLSCREEN)

	// Initialize our world
	world = make([][]int, v_info.Current_w)
	for i := range world {
		world[i] = make([]int, v_info.Current_h)
		for j := range world[i] {
			world[i][j] = OPEN
		}
	}

	// Once we're done, free screen object and quit sdl.
	defer sdl.Quit()
	defer screen.Free()

	if screen == nil {
		panic(sdl.GetError())
	}

	// Set window title
	sdl.WM_SetCaption("A* algorithm demo", "")

	// Give the screen an initially and update display
	screen.FillRect(nil, 0x33CCCC);
	screen.Flip()

	/* Draw a grid on our display */
	_, _  = drawSquare(screen)

	for true {
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {

			switch e := ev.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:
				/* Quit when escape is pressed */
				if e.Keysym.Sym == sdl.K_ESCAPE {
					return
				}
			case *sdl.MouseMotionEvent:
				if(e.State == sdl.BUTTON_LEFT) {
					drawLine(world, screen, e)
				}
			case *sdl.MouseButtonEvent:
				if e.Type == sdl.MOUSEBUTTONDOWN &&
						e.Button == sdl.BUTTON_LEFT {
					r := getRect(e)

					//fmt.Print("Click point ", r)

					if world[(r.X/20)][(r.Y/20)] == OPEN {
						fillBoxWall(world, screen, r)
					} else {
						fillBoxOpen(world, screen, r)
					}
				}
			default:
			}
		}

		// Delay for 25 milliseconds
		sdl.Delay(25)
	}

	fmt.Println("Exiting");
}

func getRect(p *sdl.MouseButtonEvent) *sdl.Rect{
	x := int(p.X)
	y := int(p.Y)

	return &sdl.Rect{
		int16((x - (x % 20)) + 1),
		int16((y - (y % 20)) + 1),
		19,
		19,
	}
}

/*
 * Draw a grid on the display and return info about the Tile
 */
func drawSquare(screen *sdl.Surface) (x, y int) {
	vid := sdl.GetVideoInfo()

	// First the vertical
	for i := 0; i < int(vid.Current_w); i += 20 {
		screen.FillRect(
			&sdl.Rect{int16(i), int16(0), 1, uint16(vid.Current_h)},
			0x000000)
		screen.UpdateRect(int32(i), 0, 1, uint32(vid.Current_h))
	}

	// Then the horizontal
	for i := 0; i < int(vid.Current_h); i += 20 {
		screen.FillRect(
			&sdl.Rect{0, int16(i), uint16(vid.Current_w), 1},
			0x000000)
			screen.UpdateRect(0, int32(i), uint32(vid.Current_w), 1)
	}

	return 20, 20
}
