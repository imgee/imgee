// Package image provides ...
package image

type Format string

var (
	PNG  Format = "png"
	BMP  Format = "bmp"
	IFF  Format = "iff"
	TIFF Format = "tiff"
	PNF  Format = "pnf"
	GIF  Format = "gif"
	JPG  Format = "jpg"
	JPEG Format = "jpeg"
	MNG  Format = "mng"
	PSD  Format = "psd"
	SAI  Format = "sai"
	UFO  Format = "ufo"
	XCF  Format = "xcf"
	PCX  Format = "pcx"
	PPM  Format = "ppm"
	WEBP Format = "webp"

	Formats = map[string]Format{
		"png":  PNG,
		"bmp":  BMP,
		"iff":  IFF,
		"tiff": TIFF,
		"pnf":  PNF,
		"gif":  GIF,
		"jpg":  JPG,
		"jpeg": JPEG,
		"mng":  MNG,
		"psd":  PSD,
		"sai":  SAI,
		"ufo":  UFO,
		"xcf":  XCF,
		"pcx":  PCX,
		"ppm":  PPM,
		"webp": WEBP,
	}
)
