// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gltext

import (
	"fmt"

	"github.com/remogatto/shaders"
	"github.com/remogatto/shapes"
)

// A Font allows rendering of text to an OpenGL context.
type Font struct {
	config         *FontConfig     // Character set for this font.
	texture        uint32          // Holds the glyph texture id.
	program        shaders.Program // Shader program
	listbase       []*shapes.Box   // Holds the first display list id.
	maxGlyphWidth  int             // Largest glyph width.
	maxGlyphHeight int             // Largest glyph height.
}

// loadFont loads the given font data. This does not deal with font scaling.
// Scaling should be handled by the independent Bitmap/Truetype loaders.
// We therefore expect the supplied image and charset to already be adjusted
// to the correct font scale.
//
// The image should hold a sprite sheet, defining the graphical layout for
// every glyph. The config describes font metadata.
func loadFont(texture Texture, config *FontConfig) (f *Font, err error) {
	f = new(Font)
	f.program = shaders.NewProgram(shapes.DefaultBoxFS, shapes.DefaultBoxVS)
	f.config = config

	// // Resize image to next power-of-two.
	// img = glh.Pow2Image(img).(*image.RGBA)
	ib := texture.Bounds()

	texWidth := float32(ib.Dx())
	texHeight := float32(ib.Dy())

	for _, glyph := range config.Glyphs {
		// Update max glyph bounds.
		if glyph.Width > f.maxGlyphWidth {
			f.maxGlyphWidth = glyph.Width
		}

		if glyph.Height > f.maxGlyphHeight {
			f.maxGlyphHeight = glyph.Height
		}

		// Quad width/height
		vw := float32(glyph.Width)
		vh := float32(glyph.Height)

		// Texture coordinate offsets.
		glyph.tx1 = float32(glyph.X) / texWidth
		glyph.ty1 = 1.0 - float32(glyph.Y)/texHeight
		glyph.tx2 = (float32(glyph.X) + vw) / texWidth
		glyph.ty2 = 1.0 - (float32(glyph.Y)+vh)/texHeight

		// Advance width (or height if we render top-to-bottom)
		// adv := float32(glyph.Advance)

		shape := shapes.NewBox(f.program, vw, vh)
		// shape.SetColor(color.White)
		shape.SetTexture(
			texture.Id(),
			[]float32{
				glyph.tx1, glyph.ty2,
				glyph.tx2, glyph.ty2,
				glyph.tx1, glyph.ty1,
				glyph.tx2, glyph.ty1,
			},
		)

		f.listbase = append(f.listbase, shape)

	}

	return
}

// Dir returns the font's rendering orientation.
func (f *Font) Dir() Direction { return f.config.Dir }

// Low returns the font's lower rune bound.
func (f *Font) Low() rune { return f.config.Low }

// High returns the font's upper rune bound.
func (f *Font) High() rune { return f.config.High }

// Glyphs returns the font's glyph descriptors.
func (f *Font) Glyphs() Charset { return f.config.Glyphs }

// // Release releases font resources.
// // A font can no longer be used for rendering after this call completes.
// func (f *Font) Release() {
// 	f.texture.Delete()
// 	gl.DeleteLists(f.listbase, len(f.config.Glyphs))
// 	f.config = nil
// }

// Metrics returns the pixel width and height for the given string.
// This takes the scale and rendering direction of the font into account.
//
// Unknown runes will be counted as having the maximum glyph bounds as
// defined by Font.GlyphBounds().
func (f *Font) Metrics(text string) (int, int) {
	if len(text) == 0 {
		return 0, 0
	}

	gw, gh := f.GlyphBounds()

	if f.config.Dir == TopToBottom {
		return gw, f.advanceSize(text)
	}

	return f.advanceSize(text), gh
}

// advanceSize computes the pixel width or height for the given single-line
// input string. This iterates over all of its runes, finds the matching
// Charset entry and adds up the Advance values.
//
// Unknown runes will be counted as having the maximum glyph bounds as
// defined by Font.GlyphBounds().
func (f *Font) advanceSize(line string) int {
	gw, gh := f.GlyphBounds()
	glyphs := f.config.Glyphs
	low := f.config.Low
	indices := []rune(line)

	var size int
	for _, r := range indices {
		r -= low

		if r >= 0 && int(r) < len(glyphs) {
			size += glyphs[r].Advance
			continue
		}

		if f.config.Dir == TopToBottom {
			size += gh
		} else {
			size += gw
		}
	}

	return size
}

// GlyphBounds returns the largest width and height for any of the glyphs
// in the font. This constitutes the largest possible bounding box
// a single glyph will have.
func (f *Font) GlyphBounds() (int, int) {
	return f.maxGlyphWidth, f.maxGlyphHeight
}

func (f *Font) Printf(format string, argv ...interface{}) (*shapes.Group, error) {
	text := shapes.NewGroup()
	str := fmt.Sprintf(format, argv...)
	indices := []rune(str)

	if len(indices) == 0 {
		return nil, nil
	}

	// Runes form display list indices.
	// For this purpose, they need to be offset by -FontConfig.Low
	low := f.config.Low
	x, y := float32(0), float32(0)
	tw, th := f.Metrics(str)
	for i := range indices {
		id := indices[i] - low
		letter := f.listbase[id].Clone()
		letter.MoveTo(x, y)
		text.Append(letter)
		adv := f.config.Glyphs[id].Advance
		x += float32(adv)
	}

	// Recalculate the center of the text group
	verts := text.Vertices()
	text.SetCenter(verts[0]+float32(tw/2), verts[1]+float32(th/2))

	return text, nil
}
