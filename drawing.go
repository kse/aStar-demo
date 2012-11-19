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
)

/*
 * Draw a grid on the display and return info about the Tile
 */
func drawGrid(screen *sdl.Surface) {
	/*{{{*/
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

	return
	/*}}}*/
}

func drawLine(w [][]Field, screen *sdl.Surface, from *Field, to *Field, color int) {
	/*{{{*/
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
		fillBox(&w[(x0 - (x0%SIZE) + 1)/SIZE][(y0 - (y0%SIZE) + 1)/SIZE], color)

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
	/*}}}*/
}

func initFillBox(screen *sdl.Surface) {
	/*{{{*/
	for i := range paint_chan {
		if i.f.T == i.c {
			continue;
		}

		i.f.T = i.c
		screen.FillRect(i.f.toRect(), uint32(i.c))
		// Do updates of screen in the main loop.
		//screen.UpdateRect(i.f.ToFourTuple());
	}
	/*}}}*/
}

func fillBox(f *Field, color int) {
	p := &Paint{f, color}
	paint_chan <- p
}

func drawMouseMotion(w [][]Field, screen *sdl.Surface, e *sdl.MouseMotionEvent) {
	/*{{{*/
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
	/*}}}*/
}
