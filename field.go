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

type Field struct{
	X int
	Y int
	T int
	left  *Field
	right *Field
	lsize int   // Size of the left tree
	rsize int   // Size of the right tree
	f  float64  // Distance from start + estimated distance to goal
	g  float64  // Distance from start
	c  bool     // Whether field is in closed set
	o  bool     // Whether field is in open set
	origin *Field
}

func (f *Field) toRect() *sdl.Rect{
	return &sdl.Rect{
		X: int16(f.X*SIZE) + 1,
		Y: int16(f.Y*SIZE) + 1,
		W: SIZE - 1,
		H: SIZE - 1,
	}
}

func (this *Field) HeapInsert(f *Field) (newRoot *Field){
	/*{{{*/
	if f == this {
		return this;
	}

	if this == nil {
		return f;
	}

	if f.f >= this.f {
		if this.lsize > this.rsize {
			if this.right == nil {
				this.right = f;
			} else {
				this.right = this.right.HeapInsert(f);
			}
			this.rsize++;
		} else {
			if this.left == nil {
				this.left = f;
			} else {
				this.left = this.left.HeapInsert(f);
			}
			this.lsize++;
		}
		newRoot = this;
	} else {
		f.right = this.right;
		f.rsize = this.rsize;

		f.left  = this.left;
		f.lsize = this.lsize;

		this.lsize = 0;
		this.rsize = 0;

		this.right = nil;
		this.left = nil

		if f.lsize > f.rsize {
			f.rsize++;
			f.right = f.right.HeapInsert(this);
		} else {
			f.lsize++;
			f.left = f.left.HeapInsert(this);
		}

		newRoot = f;
	}

	return newRoot;
	/*}}}*/
}

/*
 * Extract minimum element (the root from this heap.
 * This involves finding a new root element, and returning
 * the pointer of this
 */
func (this *Field) HeapExtractMin() (f1, newRoot *Field){
	/*{{{*/
	if this.right == nil && this.left == nil {
		// If both right and left are null, we just return ourselves
		// and a nil newRoot, because the heap is then empty
	} else if this.right == nil {
		// If our right child is null, return our left child,
		// which we know is not null.
		newRoot = this.left
	} else if this.left == nil {
		// If our left child is null, return our right child,
		// which we know is not null.
		newRoot = this.right
	} else {
		// When we're here we know that neither right nor left
		// child are nil, and it all comes down to finding the
		// minimum of the two
		if this.left.f < this.right.f {
			var newLeft *Field;
			if this == this.left {
				panic("This and left are equal");
			}
			newRoot, newLeft = this.left.HeapExtractMin();

			if newLeft != nil {
				newRoot.left  = newLeft;
				newRoot.lsize = newLeft.lsize + newLeft.rsize + 1;
			}

			newRoot.right = this.right;
			newRoot.rsize = this.rsize;
		} else {
			var newRight *Field;
			if this == this.right {
				panic("This and right are equal");
			}
			newRoot, newRight = this.right.HeapExtractMin();

			if newRight != nil {
				newRoot.right  = newRight;
				newRoot.rsize = newRight.rsize + newRight.lsize + 1;
			}

			newRoot.lsize = this.lsize;
			newRoot.left = this.left;
		}
	}

	this.lsize = 0;
	this.rsize = 0;
	this.right = nil;
	this.left = nil;

	return this, newRoot;
	/*}}}*/
}

func (f *Field) ParseRect(r *sdl.Rect, color int) {
	f.X = int(r.X)/SIZE;
	f.Y = int(r.Y)/SIZE;
	f.T = color;
}

func (f *Field) ToFourTuple() (X int32, Y int32, W uint32, H uint32){
	r := f.toRect();
	return int32(r.X), int32(r.Y), uint32(r.W), uint32(r.H);
}
