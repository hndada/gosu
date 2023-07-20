package choose

// no database. just search data linearly.

func (s Scene) drawSearchBox(screen draws.Image) {
	t := *&s.queryTypeWriter.Text()
	if t == "" {
		t = "Type for search..."
	}
	const a = "searching..."
	count := 0
	if count == 0 {
		fmt.Sprintf("found no charts", count)
	} else {
		fmt.Sprintf("%s charts found", count)
	}
	text.Draw(screen, t, scene.Face16, int(d.X), int(d.Y)+25, color.White)
}

// imported from old code.
// It was for displaying text when loading.
func loadingText(t draws.Timer) string {
	s := "Loading"
	c := int(3*t.Age() + 1)
	s += strings.Repeat(".", c)
	return s
}
