type Scene struct {
	// musicCh    chan []byte
	// bgCh       chan draws.Image
}

func (s *Scene) Update() any {
	// s.Preview.Update()
	// select {
	// case i := <-s.bgCh:
	// 	sprite := draws.NewSprite(i)
	// 	sprite.SetScaleToW(ScreenSizeX)
	// 	sprite.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
	// 	s.Background.Sprite = sprite
	// default:
	// }
	if isEnter() {
		switch s.Focus {
		case FocusChart:
			go func() {
				s.loading = true
				scene.UserSkin.Enter.Play(*s.volumeSound)
				// fs, name, err := c.Choose()
				// if err != nil {
				// 	fmt.Println(err)
				// }
				// s.Preview.Close()
				s.loading = false
			}()
		}
		return nil
	}
	// cset := s.ChartSets.Current()
	// go func() {
	// 	// i, err := draws.NewImageFromURL("https://upload.wikimedia.org/wikipedia/commons/1/1f/As08-16-2593.jpg")
	// 	i, err := draws.NewImageFromURL(cset.URLCover("cover", Large))
	// 	if err != nil {
	// 		return
	// 	}
	// 	s.Background.Sprite.Source = i
	// 	// s.bgCh <- draws.Image{Image: i}
	// }()
}

func (s *Scene) updateWeb() any {
	switch s.Focus {
	case FocusSearch:
		s.Query.Update()
	case FocusChartSet:
		fired, state := s.ChartSets.Update()
		if !fired {
			break
		}
		switch state {
		case prev:
			if s.page == 0 {
				break
			}
			s.page--
			go func() {
				s.LoadChartSetList()
				s.ChartSets.cursor = RowCount - 1
			}()
		case next:
			css := s.ChartSets
			s.page++
			go func() {
				s.LoadChartSetList()
				if len(s.ChartSets.ChartSets) == 0 {
					s.ChartSets = css
					s.page--
				}
			}()
		case stay:
		}
		cset := s.ChartSets.Current()
	}
	return nil
}

// err will be assigned to return value 'err'.
func (c Chart) Choose() (fsys fs.FS, name string, err error) {
	// const noVideo = 1
	// u := fmt.Sprintf("%s%d?n=%d", APIDownload, c.ParentSetId, noVideo)
	u := c.URLDownload()
	resp, err := http.Get(u)
	if err != nil {
		return
	}
	// fmt.Printf("download URL: %s\n", u)
	// var req *http.Request
	// var resp *http.Response
	// if runtime.GOARCH == "wasm" {
	// 	client := &http.Client{}
	// 	req, err = http.NewRequest("GET", u, nil)
	// 	if err != nil {
	// 		return
	// 	}
	// 	req.Header.Add("js.fetch:mode", "no-cors")
	// 	resp, err = client.Do(req)
	// 	if err != nil {
	// 		return
	// 	}
	// } else {
	// 	resp, err = http.Get(u)
	// 	if err != nil {
	// 		return
	// 	}
	// }
	defer resp.Body.Close()
	if !s.Preview.IsEmpty() {
		s.Preview.Close()
	}
	s.Preview, err = NewPreviewPlayer(resp.Body)
	if err != nil {
		fmt.Println("preview:", err)
		// if runtime.GOARCH != "wasm" {
		// 	return
		// }
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fsys, err = zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return
	}
	return fsys, c.OsuFile, err
}

func (c ChartSet) URLCover(kind, suffix string) string {
	return fmt.Sprintf("%s%s/%d/covers/%s%s.jpg", proxy, APIBeatmap, c.SetId, kind, suffix) // https://proxy.cors.sh/
}
func (c ChartSet) URLPreview() string {
	return fmt.Sprintf("%shttps://b.ppy.sh/preview/%d.mp3", proxy, c.SetId) // https://proxy.cors.sh/
}
func (c ChartSet) URLDownload() string {
	return fmt.Sprintf("%shttps://api.chimu.moe/v1/d/%d", proxy, c.SetId) // https://proxy.cors.sh/
}

// It goes roughly triangular number.
func Level(sr float64) int { return int(math.Pow(sr, 1.7)) }

func ChartSetList(root string) map[int]bool {
	for _, dir := range dirs {
		if dir.IsDir() || filepath.Ext(dir.Name()) == ".osz" {
			s := strings.Split(dir.Name(), " ")
			setId, err := strconv.Atoi(s[0])
		}
	}
}

type Row struct {
	// thumbCh chan draws.Image
	// cardCh  chan draws.Image

	draws.Sprite // Thumbnail
	Thumb        draws.Sprite
	Card         draws.Sprite
	Mask         draws.Sprite
	First        draws.Sprite
	Second       draws.Sprite
}

var defaultThumb = draws.Image{
	Image: ebiten.NewImage(int(RowHeight), int(RowHeight))}
var defaultCard = draws.Image{
	Image: ebiten.NewImage(int(RowWidth-RowHeight), int(RowHeight))}

func NewRow(cardURL, thumbURL, first, second string) Row {
	const thumbWidth = RowHeight // Thumbnail is a square.
	const (
		px = 5
		py = 30
	)
	r := Row{}
	r.Locate(ScreenSizeX-RowWidth, ScreenSizeY/2, draws.LeftMiddle)
	{
		s := draws.NewSprite(defaultThumb)
		s.SetSize(thumbWidth, RowHeight)
		r.Thumb = s
	}
	go func() {
		i, err := draws.NewImageFromURL(thumbURL)
		if err != nil {
			return
		}
		r.Thumb.Source = i
		// r.thumbCh <- draws.Image{Image: i}
		// close(r.thumbCh)
	}()
	{
		s := draws.NewSprite(defaultCard)
		s.SetSize(400, RowHeight)
		s.Locate(thumbWidth, 0, draws.LeftTop)
		r.Card = s
	}
	go func() {
		i, err := draws.NewImageFromURL(cardURL)
		if err != nil {
			return
		}
		r.Card.Source = i
		// r.cardCh <- draws.Image{Image: i}
		// close(r.cardCh)
	}()
	{
		s := scene.UserSkin.BoxMask
		s.SetSize(RowWidth, RowHeight)
		s.Locate(thumbWidth, 0, draws.LeftTop)
		r.Mask = s
	}
	{
		src := draws.NewText(first, scene.Face20)
		s := draws.NewSprite(src)
		s.Locate(px+thumbWidth, py, draws.LeftTop)
		r.First = s
	}
	{
		src := draws.NewText(second, scene.Face20)
		s := draws.NewSprite(src)
		s.Locate(px+thumbWidth, py-5+RowHeight/2, draws.LeftTop)
		r.Second = s
	}
	return r
}
func (r *Row) Update() {
	// select {
	// case i := <-r.thumbCh:
	// 	r.Sprite.Source = i
	// case i := <-r.cardCh:
	// 	r.Card.Source = i
	// default:
	// }
}
func (r Row) Draw(dst draws.Image) {
	r.Thumb.Position = r.Thumb.Add(r.Position)
	r.Thumb.Draw(dst, draws.Op{})
	r.Card.Position = r.Card.Add(r.Position)
	r.Card.Draw(dst, draws.Op{})
	r.Mask.Position = r.Mask.Add(r.Position)
	r.Mask.Draw(dst, draws.Op{})
	r.First.Position = r.First.Add(r.Position)
	r.First.Draw(dst, draws.Op{})
	r.Second.Position = r.Second.Add(r.Position)
	r.Second.Draw(dst, draws.Op{})
}

const (
	prev = iota - 1
	stay
	next
)

func (l *List) Update() (bool, int) {
	last := l.cursor
	fired := l.Cursor.Update()
	now := l.cursor
	if fired && now == last {
		switch now {
		case 0:
			return fired, prev
		case RowCount - 1:
			return fired, next
		}
	}
	return fired, stay
}
