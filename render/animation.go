package render

// May be useful for animation
// {
// 	fs, err := os.ReadDir("skin/bg")
// 	if err != nil {
// 		panic(err)
// 	}
// 	for _, f := range fs {
// 		if f.IsDir() || !strings.HasPrefix(f.Name(), "bg") {
// 			continue
// 		}
// 		sprite := Sprite{
// 			I: NewImage(filepath.Join("skin/bg", f.Name())),
// 		}
// 		sprite.SetFullscreen()
// 		g.DefaultBackgrounds = append(g.DefaultBackgrounds, sprite)
// 	}
// 	r := int(rand.Float64() * float64(len(g.DefaultBackgrounds)))
// 	RandomDefaultBackground = g.DefaultBackgrounds[r]
// }
