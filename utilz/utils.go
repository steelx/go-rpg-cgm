package utilz

import (
	"encoding/csv"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

/*
e.g. usage :

	heroFrames := pictures.LoadAsFrames(imgSprite, 16)
	heroSprite := pixel.NewSprite(imgSprite, heroFrames[10])
	scaledMatrix := pixel.IM.Scaled(pixel.ZV, 16)
	heroSprite.Draw(win, scaledMatrix.Moved(win.Bounds().Center()))

	*******************************OR********************************
	*** below will render sprite everytime you click on Window    ***
	*****************************************************************
	heroFrames := pictures.LoadAsFrames(imgSprite, 16)
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		tree := pixel.NewSprite(imgSprite, heroFrames[rand.Intn(len(heroFrames))])
		trees = append(trees, tree)
		matrices = append(matrices, pixel.IM.Scaled(pixel.ZV, 4).Moved(win.MousePosition()))
	}
	for i, tree := range trees {
		tree.Draw(win, matrices[i])
	}
*/

func LoadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func LoadAsFrames(imgSprite pixel.Picture, w, h float64) []pixel.Rect {
	var spriteFrames []pixel.Rect

	for y := imgSprite.Bounds().Min.Y; y < imgSprite.Bounds().Max.Y; y += h {
		for x := imgSprite.Bounds().Min.X; x < imgSprite.Bounds().Max.X; x += w {
			spriteFrames = append(spriteFrames, pixel.R(x, y, x+w, y+h))
		}
	}
	//e.g. pixel.NewSprite(imgSprite, spriteFrames[frameIndex])
	return spriteFrames
}

func LoadAsFramesFromTop(imgSprite pixel.Picture, w, h float64) []pixel.Rect {
	var spriteFrames []pixel.Rect

	minY := math.Floor(imgSprite.Bounds().Min.Y)
	maxY := math.Floor(imgSprite.Bounds().Max.Y)
	for y := maxY - h; y >= minY; y -= h {
		for x := imgSprite.Bounds().Min.X; x < imgSprite.Bounds().Max.X; x += w {
			spriteFrames = append(spriteFrames, pixel.R(x, y, x+w, y+h))
		}
	}
	//e.g. pixel.NewSprite(imgSprite, spriteFrames[frameIndex])
	return spriteFrames
}

//LoadAnimationFromCSV to set image sprite frames to good use,
// load them as set of animations
/*csv file:
Front,0,0
FrontBlink,1,1
LookUp,2,2
Left,3,7
LeftRight,4,6
LeftBlink,7,7
Walk,8,15
Run,16,23
Jump,24,26
*/
// e.g. animations = LoadAnimationFromCSV("./animations.csv", LoadAsFrames())
func LoadAnimationsFromCSV(descPath string, spriteFrames []pixel.Rect) map[string][]pixel.Rect {
	descFile, err := os.Open(descPath)
	if err != nil {
		return nil
	}
	defer descFile.Close()

	// load the animation information, Name and interval inside the spritesheet
	desc := csv.NewReader(descFile)
	var animations = make(map[string][]pixel.Rect)
	for {
		anim, err := desc.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil
		}

		name := anim[0]
		start, _ := strconv.Atoi(anim[1])
		end, _ := strconv.Atoi(anim[2])

		animations[name] = spriteFrames[start : end+1]
	}
	return animations
}

func GenerateUVs(tileWidth, tileHeight float64, texture pixel.Picture) []UV {
	// This is the table we'll fill with uvs and return.
	var uvs []UV

	textureWidth := texture.Bounds().W()
	textureHeight := texture.Bounds().H()
	width := tileWidth / textureWidth
	height := tileHeight / textureHeight
	cols := textureWidth / tileWidth
	rows := textureHeight / tileHeight

	var ux, uy float64
	var vx, vy float64 = width, height

	for rows > 0 {
		for cols > 0 {
			uvs = append(uvs, UV{ux, uy, vx, vy})
			//Advance the UVs to the next column
			ux = ux + width
			vx = vx + width
			cols -= 1
		}
		// Put the UVs back to the start of the next row
		ux = 0
		vx = width
		uy = uy + height
		vy = vy + height
		rows -= 1
	}
	return uvs
}

//LoadSprite load TMX tile image source
func LoadSprite(path string) (*pixel.Sprite, *pixel.PictureData) {
	f, err := os.Open(path)
	PanicIfErr(err)

	img, err := png.Decode(f)
	PanicIfErr(err)

	pd := pixel.PictureDataFromImage(img)
	return pixel.NewSprite(pd, pd.Bounds()), pd
}

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func RandInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}
func RandFloat(min, max float64) float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + r.Float64()*(max-min)
}
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

//Clamp restricts a number to a certain range. If a value is too high,
// it’s reduced to the maximum.
// If it’s too low, it’s increased to the minimum.
func Clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(value, max))
}

func LoadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

//HexToColor("#E53935")
func HexToColor(hex string) (c color.RGBA) {
	c.A = 0xff

	errInvalidFormat := color.RGBA{255, 255, 255, 255}

	if hex[0] != '#' {
		return errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}

		return 0
	}

	switch len(hex) {
	case 7:
		c.R = hexToByte(hex[1])<<4 + hexToByte(hex[2])
		c.G = hexToByte(hex[3])<<4 + hexToByte(hex[4])
		c.B = hexToByte(hex[5])<<4 + hexToByte(hex[6])
	case 4:
		c.R = hexToByte(hex[1]) * 17
		c.G = hexToByte(hex[2]) * 17
		c.B = hexToByte(hex[3]) * 17
	default:
		return errInvalidFormat
	}
	return
}

func GetAlpha(f float64) uint8 {
	if f >= 1 {
		return 255
	}
	return uint8(f * 256)
}

func DebugPxPoint(x, y float64, renderer pixel.Target) {
	imd := imdraw.New(nil)
	imd.Color = HexToColor("#ff00ff")
	imd.Push(pixel.V(x, y))
	imd.Circle(3, 0)
	imd.Draw(renderer)
}
