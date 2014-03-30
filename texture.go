package gltext

import "image"

type Texture interface {
	Bounds() image.Rectangle
	Id() uint32
}

type TextureUploader interface {
	UploadRGBAImage(img *image.RGBA) Texture
}
