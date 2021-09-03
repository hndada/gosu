package game

import (
	_ "image/jpeg"
	_ "image/png"
)

func (c *ChartHeader) Background() (*ebiten.Image, error) {
	dat, err := ioutil.ReadFile(c.Path(c.ImageFilename))
	if err != nil {
		return nil, err
	}
	src, _, err := image.Decode(bytes.NewReader(dat))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(src, ebiten.FilterDefault)
}

// BG 옵션은 한번 해놓으면 웬만하면 안바뀌니 저장 후 쓰는 걸로
func BackgroundOp(screen, bg image.Point) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	bx, by := float64(bg.X), float64(bg.Y)
	sx, sy := float64(screen.X), float64(screen.Y)

	rx, ry := sx/bx, sy/by
	var ratio float64 = 1
	if rx < 1 || ry < 1 { // 스크린이 그림보다 작을 경우 그림 크기 줄이기
		min := rx
		if min > ry {
			min = ry
		}
		ratio = min
		op.GeoM.Scale(ratio, ratio)
	}

	// 그림이 모니터의 중앙에 위치하게
	// x와 y 둘 중 하나는 스크린 크기와 일치; 둘 모두 크기가 스크린보다 작거나 같다
	x, y := bx*ratio, by*ratio
	op.GeoM.Translate((sx-x)/2, (sy-y)/2)
	return op
}
