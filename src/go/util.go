package main

import (
	"image"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/golang-ui/nuklear/nk"
)

func extractPSDAndPFV(filenames []string) string {
	var psd, pfv string
	for _, s := range filenames {
		switch strings.ToLower(filepath.Ext(s)) {
		case ".psd":
			psd = s
		case ".pfv":
			pfv = s
		}
	}
	if psd != "" && pfv != "" && filepath.Dir(psd) == filepath.Dir(pfv) {
		psd += "|" + filepath.Base(pfv)
	}
	return psd
}

func createTexture(nrgba *image.NRGBA) (uint32, nk.Image) {
	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, int32(nrgba.Rect.Dx()), int32(nrgba.Rect.Dy()),
		0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&nrgba.Pix[0]))
	return tex, nk.NkImageId(int32(tex))
}