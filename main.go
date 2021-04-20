package main

import (
	"fmt"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image"
	"math"
	"os"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	SIZE = 30.0
)

func loadPicture(path string) (pixel.Picture, error) {
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

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, float64(len(GMap[0])*30), float64(len(GMap)*30)),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(200, 500), basicAtlas)

	imd := imdraw.New(nil)
	imd.Color = pixel.RGBA{1, 0, 0, 0.1}
	imd.Push(pixel.V(200, 100))
	imd.Color = pixel.RGB(0, 1, 0)
	imd.Push(pixel.V(800, 100))
	imd.Color = pixel.RGB(0, 0, 1)
	imd.Push(pixel.V(500, 700))

	imd.Polygon(0)

	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			hero := MySences.FindHero(1)
			pos := GetGridByPos(win.MousePosition())
			//fmt.Fprintf(basicTxt, "Now Click  Row:%f Column:%f !\n",pos.Y,pos.X)
			getPath := FindPathIgnoreBlock([]int32{hero.location.x, int32(len(GMap)) - hero.location.y - 1}, [][]int32{{int32(pos.X), int32(len(GMap)) - int32(pos.Y) - 1}})
			allDirections := make([]Direction, 0, len(getPath)/2+1)
			movepath := make([]pixel.Vec, 0, len(allDirections))
			for {
				if len(getPath) < 2 {
					break
				}
				sliceDirection, direction := PointToDirection(getPath[0:2])
				allDirections = append(allDirections, direction)
				getPath = getPath[1:]
				movepath = append(movepath, sliceDirection)
			}
			hero.location.offset = movepath
		}
		for h := len(GMap) - 1; h >= 0; h-- {
			for w := 0; w < len(GMap[0]); w++ {
				SetMapImg([]int{w, h}, win)
			}
		}
		for _, hero := range MySences.Heros {
			hero.location.canMove = CalcMoveRange(hero)
			for _, grid := range hero.location.canMove {
				SetCanMoveImg([]int{int(grid.x), int(grid.y)}, win)
			}
			hero.MoveToNextGrid(0.1)
			DrawHeroToMap(hero, win)
		}
		imd.Draw(win)
		basicTxt.Draw(win, pixel.IM)
		win.Update()
	}
}

type Sences struct {
	mapSource map[Sourcetype]*pixel.Sprite
	Heros     []*Hero
}

var MySences *Sences

type Sourcetype int

const (
	BackGround Sourcetype = 0
	Tree       Sourcetype = 1
	Mount      Sourcetype = 2
	Move       Sourcetype = 4
)

func (s *Sences) FindHero(status int) *Hero {
	for _, hero := range s.Heros {
		if hero.status == status {
			return hero
		}
	}
	return nil
}

func init() {
	MySences = &Sences{
		mapSource: make(map[Sourcetype]*pixel.Sprite),
		Heros:     make([]*Hero, 0),
	}
	spritesheet, err := loadPicture("assest/map.png")
	if err != nil {
		panic(err)
	}
	MySences.mapSource[BackGround] = pixel.NewSprite(spritesheet, pixel.R(0, 30, 30, 60))
	MySences.mapSource[Mount] = pixel.NewSprite(spritesheet, pixel.R(0, 0, 30, 30))
	MySences.mapSource[Tree] = pixel.NewSprite(spritesheet, pixel.R(32, 32, 62, 62))

	move, err := loadPicture("assest/move.png")

	MySences.mapSource[Move] = pixel.NewSprite(move, pixel.R(0, 0, 30, 30))

	hero := &Hero{
		Hp:        100,
		Mp:        100,
		MovePoint: 3,
		location: Location{
			x:      1,
			y:      1,
			offset: nil,
		},
		ViewHero: ViewHero{Sprite: nil},
	}
	hero.status = 1

	spritesheet, _ = loadPicture("assest/hiking.png")
	hero.ViewHero.Sprite = pixel.NewSprite(spritesheet, spritesheet.Bounds())
	MySences.Heros = append(MySences.Heros, hero)
}

//画出 技能菜单
func DrawImg() {

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(100, 500), basicAtlas)
	fmt.Fprintln(basicTxt, "Hello, text!")
	fmt.Fprintln(basicTxt, "I support multiple lines!")
	fmt.Fprintf(basicTxt, "And I'm an %s, yay!", "io.Writer")
}

func (s *Sences) DrawMenu(menuList map[string]interface{}) {

}

//func(s *Sences) GetMapUnit(){
//	for _,hero := range s.Heros{
//
//	}
//}

func GetGridByPos(vec pixel.Vec) pixel.Vec {
	return pixel.Vec{
		X: math.Ceil(vec.X/SIZE) - 1,
		Y: math.Ceil(vec.Y/SIZE) - 1,
	}
}

func IsInGrid(grid [4]pixel.Vec, vec pixel.Vec) bool {
	if vec.X > grid[0].X && vec.Y < grid[0].Y &&
		vec.X < grid[1].X && vec.Y < grid[1].Y &&
		vec.X < grid[2].X && vec.Y > grid[2].Y &&
		vec.X > grid[3].X && vec.Y > grid[3].Y {
		return true
	}
	return false
}

func GridRange(grid []int32) [4]pixel.Vec {
	x := grid[0]
	y := grid[1]
	var rect [4]pixel.Vec
	rect[0] = pixel.Vec{
		X: float64(x) * 30,
		Y: float64(y) * 30,
	}
	rect[1] = pixel.Vec{
		X: float64(x+1) * 30,
		Y: float64(y) * 30,
	}
	rect[2] = pixel.Vec{
		X: float64(x + 1),
		Y: float64(y + 1),
	}
	rect[3] = pixel.Vec{
		X: float64(x),
		Y: float64(y + 1),
	}
	return rect
}

type OBJECT int

const (
	GRASS OBJECT = 0
	TREE  OBJECT = 1
	MOUNT OBJECT = 2
)

var GMap = [][]int32{
	{1, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{1, 1, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{2, 2, 2, 2, 2, 1, 2, 2, 1, 2, 2, 2, 2, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 0, 1, 0, 2, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 0, 1, 0, 2, 2, 2, 2, 2, 2, 2, 1, 2, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 2, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 2, 1, 0, 0, 1, 0, 2, 0, 0, 0, 1, 2, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 0, 1, 0, 2, 2, 2, 0, 0, 0, 2, 1, 2, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 0, 2, 2, 0, 1, 0, 2, 2, 2, 0, 1, 2, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 0, 2, 0, 0, 1, 0, 0, 2, 2, 0, 1, 2, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 0, 2, 0, 0, 1, 0, 0, 2, 0, 0, 1, 2, 2, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 0, 2, 0, 0, 1, 2, 0, 2, 0, 2, 1, 0, 2, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 0, 2, 1, 1, 1, 2, 0, 2, 0, 2, 2, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 1, 2, 2, 2, 1, 0, 0, 1, 0, 0, 2, 0, 0, 2, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 2, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{1, 1, 0, 1, 0, 1, 0, 0, 1, 0, 0, 2, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 2, 0, 0, 1, 0, 0, 1, 2, 2, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 2, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 2, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
}

func main() {
	pixelgl.Run(run)
}
