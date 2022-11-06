package draws

import "github.com/hajimehoshi/ebiten/v2"

type Animation []Sprite

func NewAnimation(path string) Animation {
	return NewAnimationFromImages(NewImages(path))
}
func NewAnimationFromImages(images []*ebiten.Image) (a Animation) {
	a = make(Animation, len(images))
	for i, image := range images {
		a[i] = NewSpriteFromImage(image)
	}
	return
}

// func NewAnimation0(path string) (a Animation) {
// 	const ext = ".png"
// 	one := Animation{NewSprite(path + ext)}
// 	dir, err := os.Open(path)
// 	if err != nil {
// 		return one
// 	}
// 	defer dir.Close()
// 	fs, err := dir.ReadDir(-1)
// 	if err != nil {
// 		return one
// 	}

// 	nums := make([]int, 0, len(fs))
// 	for _, f := range fs {
// 		if f.IsDir() {
// 			continue
// 		}
// 		num := strings.TrimSuffix(f.Name(), ext)
// 		if num, err := strconv.Atoi(num); err == nil {
// 			nums = append(nums, num)
// 		}
// 	}
// 	sort.Ints(nums)
// 	for _, num := range nums {
// 		path := filepath.Join(path, fmt.Sprintf("%d.png", num))
// 		a = append(a, NewSprite(path))
// 	}
// 	return
// }
