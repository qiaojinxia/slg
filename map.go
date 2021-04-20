package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func SetMapImg(location []int, window *pixelgl.Window) {
	x := location[0]
	y := location[1]
	mapType := GMap[y][x]
	switch mapType {
	case 0:
		MySences.mapSource[Sourcetype(0)].Draw(window, pixel.IM.Moved(pixel.Vec{
			X: float64(x)*30 + 13,
			Y: float64(len(GMap)-y-1)*30 + 13,
		}))
	case 1:
		MySences.mapSource[Sourcetype(1)].Draw(window, pixel.IM.Moved(pixel.Vec{
			X: float64(x)*30 + 13,
			Y: float64(len(GMap)-y-1)*30 + 13,
		}))
	case 2:
		MySences.mapSource[Sourcetype(2)].Draw(window, pixel.IM.Moved(pixel.Vec{
			X: float64(x)*30 + 13,
			Y: float64(len(GMap)-y-1)*30 + 13,
		}))
	}

}

func SetCanMoveImg(location []int, window *pixelgl.Window) {
	x := location[0]
	y := location[1]
	MySences.mapSource[Sourcetype(4)].Draw(window, pixel.IM.Moved(pixel.Vec{
		X: float64(x)*30 + 13,
		Y: float64(y)*30 + 13,
	}))

}

func DrawHeroToMap(hero *Hero, window *pixelgl.Window) {
	pos_x := ((float64(hero.location.x) + hero.location.offset_x) * SIZE) + (SIZE / 2)
	pos_y := ((float64(hero.location.y) + hero.location.offset_y) * SIZE) + (SIZE / 2)
	hero.Draw(window, pixel.IM.Scaled(pixel.ZV, 0.05).Moved(pixel.Vec{pos_x, pos_y}))
}

type MoveS struct {
	x         int32
	y         int32
	movepoint int
}

func OverMap(vec pixel.Vec) bool {
	height := len(GMap)
	width := len(GMap[0])
	if vec.X < 0 || vec.Y < 0 || int(vec.Y) > height || int(vec.X) > width {
		return true
	}
	return false
}

func CanMove(vec pixel.Vec) bool {
	if int(vec.Y) > len(GMap)-1 || int(vec.X) > len(GMap[0])-1 || vec.X < 0 || vec.Y < 0 {
		return false
	}
	mapRow := GMap[len(GMap)-1-int(vec.Y)]
	column := mapRow[int(vec.X)]
	switch OBJECT(column) {
	case GRASS:
		return true
	case MOUNT:
		return false
	case TREE:
		return true
	}
	return false
}

func CalcMoveRange(hero *Hero) map[string]*MoveS {
	openlist := make([]*MoveS, 0)
	list := &MoveS{
		x:         hero.location.x,
		y:         hero.location.y,
		movepoint: hero.MovePoint,
	}
	openlist = append(openlist, list)
	openlist_recorder := make(map[string]struct{})
	closelist := make(map[string]*MoveS)

	//将 初始点 加入关闭列表
	for {
		if len(openlist) < 1 {
			break
		}
		//遍历开放列表的各个点
		centerpoint := openlist[0]
		if centerpoint.movepoint < 1 {
			closelist[fmt.Sprintf("%d%d", centerpoint.x, centerpoint.y)] = centerpoint
			continue
		}
		for _, direction := range GDir {
			new_x := centerpoint.x + direction[0]
			new_y := centerpoint.y + direction[1]

			point := pixel.Vec{float64(new_x), float64(new_y)}
			//超出地图处理
			if OverMap(point) || !CanMove(point) {
				continue
			}
			//计算不同地图的行动点消耗
			n_point := &MoveS{
				x:         new_x,
				y:         new_y,
				movepoint: centerpoint.movepoint - 1,
			}
			_, exists := closelist[fmt.Sprintf("%d%d", centerpoint.x, centerpoint.y)]
			if exists {
				continue
			}
			if n_point.movepoint <= 0 {
				closelist[fmt.Sprintf("%d%d", n_point.x, n_point.y)] = n_point
			} else {
				if _, exists := openlist_recorder[fmt.Sprintf("%d%d", n_point.x, n_point.y)]; !exists {
					openlist = append(openlist, n_point)
					openlist_recorder[fmt.Sprintf("%d%d", n_point.x, n_point.y)] = struct{}{}
				}
			}

		}
		closelist[fmt.Sprintf("%d%d", centerpoint.x, centerpoint.y)] = centerpoint
		openlist = openlist[1:]
	}
	return closelist
	//从当前地点便利
}

//在屏幕上画出菜单
func DrawMenu(window *pixelgl.Window) {
	//	imd := imdraw.New(nil)
	//
	//}
	//	img, err := png.Decode(f)
	//	if err != nil {
	//		panic(err)
	//	pd := pixel.PictureDataFromImage(img)
	return
}
