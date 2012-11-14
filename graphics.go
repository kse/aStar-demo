package main

import (
	"github.com/banthar/Go-SDL/sdl"
	//pqueue "github.com/nu7hatch/gopqueue"
	"fmt"
)

const (
	OPEN  = 0x33CCCC
	WALL  = 0xFF4040
	START = 0xFF7373
	GOAL  = 0x679B00
	LINE  = 0xA60000
	// Size of rectangles	
	SIZE  = 5
)

type Field struct{
	X int
	Y int
	T int
	p int
	right *Field
	left  *Field
	lsize int
	rsize int
	prio  float32
}

func (this *Field) HeapInsert(f *Field) {
	if f == this {
		return;
	}

	if f.prio > this.prio {
		if this.lsize > this.rsize {
			if this.right == nil {
				this.right = f;
			} else {
				this.right.HeapInsert(f);
			}
			this.rsize++;
		} else {
			if this.left == nil {
				this.left = f;
			} else {
				this.left.HeapInsert(f);
			}
			this.lsize++;
		}
	} else {
		f.right = this.right;
		f.rsize = this.rsize;

		f.left = this.left;
		f.lsize = this.lsize;

		this.lsize = 0;
		this.rsize = 0;

		if f.lsize > f.lsize {
			f.right.HeapInsert(this);
		} else {
			f.left.HeapInsert(this);
		}
	}
}

func (this *Field) HeapExtractMin() (f1, f2 *Field){
	var newRoot *Field = nil;

	if this.right == nil && this.left == nil {
	} else if this.right == nil {
		newRoot = this.left
	} else if this.left == nil {
		newRoot = this.right
	} else {
		if this.left.prio < this.right.prio {
			newRoot = this.left;
			newRoot, newLeft := this.left.HeapExtractMin();

			newRoot.lsize = newLeft.lsize + newLeft.rsize + 1;
			newRoot.left  = newLeft;

			newRoot.rsize = this.rsize;
			newRoot.right = this.right;
		} else {
			newRoot = this.right;
			newRoot, newRight := this.right.HeapExtractMin();

			newRoot.rsize = newRight.rsize + newRight.lsize + 1;
			newRoot.right  = newRight;

			newRoot.lsize = this.lsize;
			newRoot.left = this.left;
		}
	}

	this.lsize = 0;
	this.rsize = 0;
	this.right = nil;
	this.left = nil;

	return this, newRoot;
}

func (f *Field) ParseRect(r *sdl.Rect, color int) {
	f.X = int(r.X)/SIZE;
	f.Y = int(r.Y)/SIZE;
	f.T = color;
}

func (f *Field) toRect() *sdl.Rect{
	return &sdl.Rect{
		X: int16(f.X*SIZE) + 1,
		Y: int16(f.Y*SIZE) + 1,
		W: SIZE - 1,
		H: SIZE - 1,
	}
}

func (f *Field) ToFourTuple() (X int32, Y int32, W uint32, H uint32){
	r := f.toRect();
	//fmt.Println("Return fourtuple:", int32(r.X), int32(r.Y), uint32(r.W), uint32(r.H));
	return int32(r.X), int32(r.Y), uint32(r.W), uint32(r.H);
}

func (f *Field) GetNeighbours(world [][]Field, ch chan<- *Field) {
}

/*
 * The star of the show!
*/
func aStar(w [][]Field, screen *sdl.Surface, start *Field, goal *Field) {
	drawLine(w, screen, start, goal, LINE)
	fillBox(screen, &w[start.X][start.Y], START);
	fillBox(screen, &w[goal.X][goal.Y], GOAL);
}

func main() {
	// Contains our world, which is simply an array of types
	var world [][]Field
	var start *Field;
	var goal *Field;

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
	world = make([][]Field, v_info.Current_w)
	for i := range world {
		world[i] = make([]Field, v_info.Current_h)
		for j := range world[i] {
			world[i][j].X = i;
			world[i][j].Y = j;
			world[i][j].T = OPEN;
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
				if e.State == sdl.BUTTON_LEFT || e.State == sdl.BUTTON_WHEELUP {
					drawMouseMotion(world, screen, e)
				}
			case *sdl.MouseButtonEvent:
				if e.Type == sdl.MOUSEBUTTONDOWN {
					state := sdl.GetKeyState()
					if e.Button == sdl.BUTTON_LEFT {
						r := getRect(e)

						if state[sdl.K_s] == 1 {
							// Left mouse button with s, set new start point
							if start == nil {
								start = new(Field)
								start.ParseRect(r, OPEN)

								fillBox(screen, start, START);
							} else {
								fillBox(screen, start, OPEN);

								start.ParseRect(r, OPEN);
								fillBox(screen, start, START);
							}

							if start != nil && goal != nil {
								aStar(world, screen, start, goal)
							}
						} else if state[sdl.K_g] == 1 {
							// Left mouse button with g, set new goal point
							if goal == nil {
								goal = new(Field)
								goal.ParseRect(r, START)

								fillBox(screen, goal, GOAL);
							} else {
								fillBox(screen, goal, OPEN);
								goal.ParseRect(r, OPEN);
								fillBox(screen, goal, GOAL);
							}

							if start != nil && goal != nil {
								aStar(world, screen, start, goal)
							}
						} else {
							// No relevant modifiers were pressed, color the field.
							var f *Field = &world[(r.X/SIZE)][(r.Y/SIZE)];
							//fmt.Println("Click on", f);

							if f.T == OPEN {
								fillBox(screen, f, WALL)
							} else {
								fillBox(screen, f, OPEN)
							}
						}
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
		int16((x - (x % SIZE)) + 1),
		int16((y - (y % SIZE)) + 1),
		SIZE - 1,
		SIZE - 1,
	}
}

/*
 * Draw a grid on the display and return info about the Tile
 */
func drawSquare(screen *sdl.Surface) (x, y int) {
	vid := sdl.GetVideoInfo()

	// First the vertical
	for i := 0; i < int(vid.Current_w); i += SIZE {
		screen.FillRect(
			&sdl.Rect{int16(i), int16(0), 1, uint16(vid.Current_h)},
			0x000000)
		screen.UpdateRect(int32(i), 0, 1, uint32(vid.Current_h))
	}

	// Then the horizontal
	for i := 0; i < int(vid.Current_h); i += SIZE {
		screen.FillRect(
			&sdl.Rect{0, int16(i), uint16(vid.Current_w), 1},
			0x000000)
			screen.UpdateRect(0, int32(i), uint32(vid.Current_w), 1)
	}

	return SIZE, SIZE
}

func drawLine(w [][]Field, screen *sdl.Surface, from *Field, to *Field, color int) {
	var x1, y1 int = to.X*SIZE, to.Y*SIZE;
	var x0, y0 int = from.X*SIZE, from.Y*SIZE;
	var sx, sy int
	var err, e2 int
	//var f *Field = new(Field)

	dx := abs(x1 - x0);
	dy := abs(y1 - y0);

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

	if dx > dy {
		err = dx/2
	} else {
		err = -dy/2
	}

	for true {
		fillBox(screen, &w[(x0 - (x0%SIZE) + 1)/SIZE][(y0 - (y0%SIZE) + 1)/SIZE], color)

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
}

//func fillBox(w [][]int, screen *sdl.Surface, r *sdl.Rect, color int) {
func fillBox(screen *sdl.Surface, f *Field, color int) {
	if f.T == color {
		return
	}

	f.T = color
	screen.FillRect(f.toRect(), uint32(color))
	screen.UpdateRect(f.ToFourTuple());
}

func drawMouseMotion(w [][]Field, screen *sdl.Surface, e *sdl.MouseMotionEvent) {
	var color int;
	if e.State == sdl.BUTTON_LEFT {
		color = WALL
	} else {
		color = OPEN
	}

	drawLine(w, screen,
		&Field{
			X: int(e.X) / SIZE,
			Y: int(e.Y) / SIZE,
		},
		&Field{
			X: (int(e.X) - int(e.Xrel))/SIZE,
			Y: (int(e.Y) - int(e.Yrel))/SIZE,
		},
		color)
}

func abs(v int) int{
	if v < 0 {
		return v * -1
	}

	return v
}
