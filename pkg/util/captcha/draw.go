package captcha

import (
	"bytes"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math/rand"
	"os"
	"time"
)

// 使用私有的随机生成器
var irand = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	defaultChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	adjust       = 12
	fontSize     = adjust * 2
)

type Drawer struct {
	font   *truetype.Font
	bg     image.Image
	chars  []rune
	length int
}

// NewDrawer font为ttf格式必传，background为jpg格式默认纯灰，chars支持中文默认26个英文大写字母
func NewDrawer(font, background, chars string) *Drawer {
	b, err := os.ReadFile(font)
	if err != nil {
		panic(err)
	}
	f, err := freetype.ParseFont(b)
	if err != nil {
		panic(err)
	}
	d := &Drawer{font: f}
	if background != "" {
		if source, err := os.Open(background); err == nil {
			bg, _ := jpeg.Decode(source)
			source.Close()
			d.bg = bg
		}
	}
	if chars == "" {
		d.chars = []rune(defaultChars)
	} else {
		d.chars = []rune(chars)
	}
	d.length = len(d.chars)
	return d
}

// Generate 绘制验证码
func (d *Drawer) Generate(n int) (string, []byte) {
	height := fontSize * 2
	width := fontSize*(n+1) - adjust/2
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	if d.bg == nil {
		draw.Draw(canvas, canvas.Bounds(), image.NewUniform(color.Gray{Y: uint8(irand.Intn(256))}),
			image.Pt(0, 0), draw.Src)
	} else {
		x, y := 0, 0
		if bgx := d.bg.Bounds().Dx(); bgx > width {
			x = irand.Intn(bgx - width)
		}
		if bgy := d.bg.Bounds().Dy(); bgy > height {
			y = irand.Intn(bgy - height)
		}
		draw.Draw(canvas, canvas.Bounds(), d.bg, image.Pt(x, y), draw.Src)
	}

	c := freetype.NewContext()
	c.SetDst(canvas)
	c.SetClip(canvas.Bounds())
	c.SetFont(d.font)
	c.SetFontSize(fontSize)

	x := adjust / 2
	code := make([]rune, n)
	for i := range code {
		char := d.chars[irand.Intn(d.length)]
		code[i] = char
		pt := freetype.Pt(x+irand.Intn(adjust), fontSize+irand.Intn(adjust))
		c.SetSrc(randColor())
		c.DrawString(string(char), pt) // nolint
		x += fontSize
	}

	buf := bytes.NewBuffer(nil)
	_ = jpeg.Encode(buf, canvas, nil)
	return string(code), buf.Bytes()
}

func randColor() *image.Uniform {
	return image.NewUniform(color.RGBA{
		R: uint8(irand.Intn(256)),
		G: uint8(irand.Intn(256)),
		B: uint8(irand.Intn(256)),
		A: 255,
	})
}
