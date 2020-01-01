package main

import (
	"encoding/csv"
	"fmt"
	"github.com/faiface/pixel"
	"image"
	"image/png"
	_ "image/png"
	"io"
	"os"
	"strconv"
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

	// load the animation information, name and interval inside the spritesheet
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
	panicIfErr(err)

	img, err := png.Decode(f)
	panicIfErr(err)

	pd := pixel.PictureDataFromImage(img)
	return pixel.NewSprite(pd, pd.Bounds()), pd
}

func panicIfErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
