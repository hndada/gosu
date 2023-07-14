func Mode(path string) (int, int) {
	const modeError = -1

	f, err := os.Open(path)
	if err != nil {
		return modeError, 0
	}
	defer f.Close()
	var (
		mode     int
		keyCount int
	)
	scanner := bufio.NewScanner(f) // Splits on newlines by default.
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "Mode: ") {
			vs := strings.Split(scanner.Text(), "Mode: ")
			if len(vs) != 2 {
				return ModeDefault, 0 // Blank goes default mode.
			}
			v, err := strconv.Atoi(vs[1])
			if err != nil {
				return modeError, 0
			}
			mode = v
		}
		if strings.HasPrefix(scanner.Text(), "CircleSize:") {
			vs := strings.Split(scanner.Text(), "CircleSize:")
			if len(vs) != 2 {
				return mode, 0
			}
			v, err := strconv.Atoi(vs[1])
			if err != nil {
				return mode, 0
			}
			keyCount = v
			return mode, keyCount
		}
	}
	if err := scanner.Err(); err != nil {
		return modeError, 0
	}
	return ModeDefault, 0
}
