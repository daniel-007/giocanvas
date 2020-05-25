// Package giocanvas is a 2D canvas API built on gio
package giocanvas

import (
	"image"
	"image/color"
	_ "image/gif" // needed by image
	_ "image/jpeg"
	_ "image/png"
	"os"

	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// Canvas defines the Gio canvas
type Canvas struct {
	Width, Height float32
	TextColor     color.RGBA
	Context       *layout.Context
}

// NewCanvas initializes a Canvas
func NewCanvas(width, height float32) *Canvas {
	gofont.Register()
	canvas := new(Canvas)
	canvas.Width = width
	canvas.Height = height
	canvas.TextColor = color.RGBA{0, 0, 0, 255}
	canvas.Context = new(layout.Context)
	return canvas
}

// textops places text
func (c *Canvas) textops(x, y, size float32, alignment text.Alignment, s string, color color.RGBA) {
	offset := x
	th := material.NewTheme()
	switch alignment {
	case text.End:
		offset = x - c.Width
	case text.Middle:
		offset = x - c.Width/2
	}
	var stack op.StackOp
	stack.Push(c.Context.Ops)
	op.TransformOp{}.Offset(f32.Point{X: offset, Y: y}).Add(c.Context.Ops)
	l := material.Label(th, unit.Dp(size), s)
	l.Color = color
	l.Alignment = alignment
	l.Layout(c.Context)
	stack.Pop()
}

// Text places text at (x,y)
func (c *Canvas) Text(x, y, size float32, s string, color color.RGBA) {
	c.textops(x, y, size, text.Start, s, color)
}

// TextMid places text centered at (x,y)
func (c *Canvas) TextMid(x, y, size float32, s string, color color.RGBA) {
	c.textops(x, y, size, text.Middle, s, color)
}

// TextEnd places text aligned to the end
func (c *Canvas) TextEnd(x, y, size float32, s string, color color.RGBA) {
	c.textops(x, y, size, text.End, s, color)
}

// Rect makes a filled Rectangle; left corner at (x, y), with dimensions (w,h)
func (c *Canvas) Rect(x, y, w, h float32, color color.RGBA) {
	ops := c.Context.Ops
	r := f32.Rect(x, y+h, x+w, y)
	paint.ColorOp{Color: color}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// CenterRect makes a filled rectangle centered at (x, y), with dimensions (w,h)
func (c *Canvas) CenterRect(x, y, w, h float32, color color.RGBA) {
	ops := c.Context.Ops
	r := f32.Rect(x-(w/2), y+(h/2), x+(w/2), y-(h/2))
	paint.ColorOp{Color: color}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// VLine makes a vertical line beginning at (x,y) with dimension (w, h)
// half of the width is left of x, the other half is the to right of x
func (c *Canvas) VLine(x, y, w, h float32, color color.RGBA) {
	c.Rect(x-(w/2), y, w, h, color)
}

// HLine makes a horizontal line starting at (x, y), with dimensions (w, h)
// half of the height is above y, the other below
func (c *Canvas) HLine(x, y, w, h float32, color color.RGBA) {
	c.Rect(x, y-(h/2), w, h, color)
}

// Grid uses horizontal and vertical lines to make a grid
func (c *Canvas) Grid(width, height, size, interval float32, color color.RGBA) {
	var x, y float32
	for y = 0; y <= height; y += height / interval {
		c.HLine(0, y, width, size, color)
	}
	for x = 0; x <= width; x += width / interval {
		c.VLine(x, 0, size, height, color)
	}
}

// CenterImage places a named image centered at (x, y)
// scaled using the specified dimensions (w, h)
func (c *Canvas) CenterImage(name string, x, y float32, w, h int, scale float32) {
	r, err := os.Open(name)
	if err != nil {
		return
	}
	im, _, err := image.Decode(r)
	if err != nil {
		return
	}
	// compute scaled image dimensions
	// if w and h are zero, use the natural dimensions
	sc := scale / 100
	imw := float32(w) * sc
	imh := float32(h) * sc
	if w == 0 && h == 0 {
		b := im.Bounds()
		imw = float32(b.Max.X) * sc
		imh = float32(b.Max.Y) * sc
	}
	// center the image
	x = x - (imw / 2)
	y = y - (imh / 2)
	ops := c.Context.Ops
	paint.NewImageOp(im).Add(ops)
	paint.PaintOp{Rect: f32.Rect(x, y, x+imw, y+imh)}.Add(ops)
}
