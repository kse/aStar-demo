/** 
  * Copyright (c) 2012, Kasper Sacharias Eenberg
  * All rights reserved.
  * 
  * Redistribution and use in source and binary forms, with or without
  * modification, are permitted provided that the following conditions are met:
  * 
  * - Redistributions of source code must retain the above copyright notice,
  *   this list of conditions and the following disclaimer.
  * 
  * - Redistributions in binary form must reproduce the above copyright notice,
  *   this list of conditions and the following disclaimer in the documentation
  *   and/or other materials provided with the distribution.
  * 
  * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
  * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
  * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
  * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
  * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
  * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
  * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
  * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
  * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
  * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
  * POSSIBILITY OF SUCH DAMAGE.
  *
  **/

package main

import (
	"github.com/banthar/Go-SDL/sdl"
	"fmt"
	"math"
)

const (
	OPEN  = 0x9FEE00
	WALL  = 0xA60000

	START = 0x008500
	GOAL  = 0x95004B

	PATH = 0xFF0000

	OPENSET   = 0xE6399B
	CLOSEDSET = 0xCD0074

	// Size of rectangles	
	SIZE  = 5
)

type Paint struct {
	f *Field
	c int
}

var paint_chan chan *Paint;
var field_chan chan *Field;
var read_field chan *Field;

// Contains our world, which is simply an array of types
var world [][]Field
var start *Field;
var goal *Field;
var rows, columns int32;

func isPassable(f *Field) bool{
	switch(f.T) {
	case OPENSET:
		return true;
	case CLOSEDSET:
		return true;
	case OPEN:
		return true;
	}
	return false;
}

func GetNeighbours(w [][]Field, ch chan<- *Field) {
	/*{{{*/
	for f := range field_chan {
		lx := (len(w) - 1)
		ly := (len(w[0]) - 1)
		//fmt.Println("Width is", lx, "height is", ly);

		if f.X < lx {
			ch <- &w[f.X + 1][f.Y];

			if f.Y > 0 {
				ch <- &w[f.X][f.Y-1];
				if isPassable(&w[f.X + 1][f.Y]) || isPassable(&w[f.X][f.Y-1]) {
					ch <- &w[f.X + 1][f.Y-1];
				}
			}
		}

		if f.Y > 0 {
			ch <- &w[f.X][f.Y-1];
		}

		if f.Y < ly {
			ch <- &w[f.X][f.Y + 1]

			if f.X > 0 {
				if isPassable(&w[f.X][f.Y + 1]) || isPassable(&w[f.X - 1][f.Y]) {
					ch <- &w[f.X - 1][f.Y + 1];
				}
			}
		}

		if f.X > 0 {
			ch <- &w[f.X - 1][f.Y];
		}

		if f.Y > 0 && f.X > 0 {
			if isPassable(&w[f.X - 1][f.Y]) || isPassable(&w[f.X][f.Y - 1]) {
				ch <- &w[f.X - 1][f.Y - 1];
			}
		}

		if f.Y < ly && f.X < lx {
			if isPassable(&w[f.X + 1][f.Y]) || isPassable(&w[f.X][f.Y + 1]) {
				ch <- &w[f.X + 1][f.Y + 1];
			}
		}

		ch <- nil;
	}
	/*}}}*/
}

func reset(t func(int) bool) {
	for i:= range world {
		for j := range world[i] {
			w := &world[i][j];
			if t(w.T) {
				fillBox(w, OPEN);
			}
			w.c = false;
			w.o = false;
			w.f = 0.0
			w.g = 0
			w.left = nil
			w.right = nil
			w.lsize = 0
			w.rsize = 0
			w.origin = nil;
		}
	}
}

func resetAllPaths() {
	f := func(t int) (b bool) {
		switch (t) {
		case OPEN:
			return false;
		case WALL:
			return false;
		}
		return true
	}
	reset(f);
	start = nil;
	goal = nil;
}

func resetPaths() {
	f := func(t int) (b bool) {
		switch (t) {
		case OPEN:
			return false;
		case WALL:
			return false;
		case GOAL:
			return false;
		case START:
			return false;
		}
		return true
	}
	reset(f);
}

func resetComplete() {
	f := func(t int) (b bool) {
		switch (t) {
		case OPEN:
			return false;
		}
		return true
	}
	reset(f);
	start = nil;
	goal = nil;
}

/*
 * The star of the show!
*/
func aStar(w [][]Field, screen *sdl.Surface) {
	/*{{{*/
	var q, min *Field;

	q = q.HeapInsert(start);
	for q != nil {
		min, q = q.HeapExtractMin();

		if min.X == goal.X && min.Y == goal.Y{
			for min.origin != nil {
				min = min.origin;
				fillBox(min, PATH);
			}

			fillBox(goal, GOAL)
			fillBox(start, START)

			//goal = nil;
			//start = nil;
			return
		}

		min.c = true
		if min.T != START {
			fillBox(min, CLOSEDSET)
		}

		field_chan <- min;
		for f := range read_field {
			var tg float64;
			if f == nil {
				break;
			}

			if f.c {
				continue;
			}

			if f.T == WALL {
				continue;
			}

			if abs(min.X - f.X) == 1 && abs(min.Y - f.Y) == 1 {
				tg = min.g + 1.4;
			} else {
				tg = min.g + 1;
			}

			if f.o == false || tg < f.g {
				f.origin = min;
				f.g = tg;
				if f.f == 0 {
					f.f = float64(tg) + math.Hypot(float64(f.X - goal.X), float64(f.Y - goal.Y));
				}
				if !f.o {
					f.o = true;
					q = q.HeapInsert(f);
					fillBox(f, OPENSET)
				}
			}
		}
	}

	return;
	/*}}}*/
}

func main() {
	if sdl.Init(sdl.INIT_VIDEO) != 0 {
		panic(sdl.GetError())
	}

	v_info := sdl.GetVideoInfo()

	var screen = sdl.SetVideoMode(
		int(v_info.Current_w),
		int(v_info.Current_h),
		32,
		sdl.HWSURFACE | sdl.DOUBLEBUF)

	rows = v_info.Current_w / SIZE;
	columns = v_info.Current_h / SIZE;

	if v_info.Current_w % SIZE != 0 {
		rows += 1;
	}

	if v_info.Current_h % SIZE != 0 {
		columns += 1;
	}

	// Initialize our world
	world = make([][]Field, rows)
	for i := range world {
		world[i] = make([]Field, columns)
		for j := range world[i] {
			world[i][j].X     = i;
			world[i][j].Y     = j;
			world[i][j].T     = OPEN;
			world[i][j].lsize = 0;
			world[i][j].rsize = 0;
			world[i][j].o     = false;
			world[i][j].c     = false;
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
	screen.FillRect(nil, OPEN);
	screen.Flip()

	/* Draw the grid on our display */
	drawGrid(screen)

	/* Create the different channels we need */
	paint_chan = make(chan *Paint, 2000);
	read_field = make(chan *Field);
	field_chan = make(chan *Field);

	// Fillbox runs asynchronously, start it here
	go initFillBox(screen);

	// Getneighbours process runs all the time, starts here.
	go GetNeighbours(world, read_field);

	for true {
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {

			switch e := ev.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:
				if e.Keysym.Sym == sdl.K_ESCAPE {
					/* Quit when escape is pressed */
					return
				} else if e.Keysym.Sym == sdl.K_r &&
					(e.Keysym.Mod & sdl.KMOD_LCTRL) != 0 &&
					(e.Keysym.Mod & sdl.KMOD_LSHIFT) != 0 {
					resetComplete();
				} else if e.Keysym.Sym == sdl.K_r && (e.Keysym.Mod & sdl.KMOD_LCTRL) != 0 {
					resetAllPaths();
				} else if e.Keysym.Sym == sdl.K_r {
					resetPaths();
				} else if e.Keysym.Sym == sdl.K_RETURN {
					/* If RETURN is pressed, run pathfinding */
					if start != nil && goal != nil {
						//resetPaths();
						go aStar(world, screen)
					}
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
								start = &world[int(r.X)/SIZE][int(r.Y)/SIZE];

								fillBox(start, START);
							} else {
								fillBox(start, OPEN);

								start = &world[int(r.X)/SIZE][int(r.Y)/SIZE];
								fillBox(start, START);
							}
						} else if state[sdl.K_g] == 1 {
							// Left mouse button with g, set new goal point
							if goal == nil {
								goal = &world[int(r.X)/SIZE][int(r.Y)/SIZE];

								fillBox(goal, GOAL);
							} else {
								fillBox(goal, OPEN);

								goal = &world[int(r.X)/SIZE][int(r.Y)/SIZE];
								fillBox(goal, GOAL);
							}
						} else {
							// No relevant modifiers were pressed, color the field.
							var f *Field = &world[(r.X/SIZE)][(r.Y/SIZE)];
							//fmt.Println("Click on", f);

							if f.T == OPEN {
								fillBox(f, WALL)
							} else {
								fillBox(f, OPEN)
							}
						}
					}
				}
			default:
			}
		}

		// Delay for 15 milliseconds
		sdl.Delay(30)
		screen.Flip();
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

func abs(v int) int{
	if v < 0 {
		return v * -1
	}

	return v
}
