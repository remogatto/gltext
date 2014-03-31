## gltext

**Note**: This package is experimental and subject to change.
Use at your own discretion.

The gltext package offers a simple set of text rendering utilities for
OpenGL ES 2.0 programs. It deals with TrueType and Bitmap (raster)
fonts (the last one is not yet implemented on this fork). 

The package supports the full set of unicode characters, provided the
loaded font does as well.


### TODO

* Have a look at Valve's 'Signed Distance Field` techniques to render
  sharp font textures are different zoom levels.

  * [SIGGRAPH2007_AlphaTestedMagnification.pdf](http://www.valvesoftware.com/publications/2007/SIGGRAPH2007_AlphaTestedMagnification.pdf)
  * [Youtube video](http://www.youtube.com/watch?v=CGZRHJvJYIg)
  
  More links to info in the youtube video description. An alternative
  might be a port of [GLyphy](http://code.google.com/p/glyphy/)


### Known bugs

* Determining the height of truetype glyphs is not entirely accurate.
  It is unclear at this point how to get to this information reliably.
  Specifically the parts in `LoadTruetype` at truetype.go#L54+.  The
  vertical glyph bounds computed by freetype-go are not correct for
  certain fonts. Right now we manually offset the value by added `4`
  to the height. This is an unreliable hack and should be fixed.
* `freetype-go` does not expose `AdvanceHeight` for vertically
  rendered fonts.  This may mean that the Advance size for
  top-to-bottom fonts is incorrect.


### Dependencies

	go get code.google.com/p/freetype-go
    

### Usage

    go get github.com/go-gl/gltext

Refer to [go-gl/examples/gltext][ex] for usage examples.

[ex]: https://github.com/go-gl/examples/tree/master/gltext


### License

Copyright 2012 The go-gl Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.

