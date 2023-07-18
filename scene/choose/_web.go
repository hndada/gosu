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