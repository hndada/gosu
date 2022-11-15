// Temporary function.
func SetKeySettings(props []mode.Prop) {
	data, err := os.ReadFile("keys.txt")
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	for _, line := range strings.Split(string(data), "\n") {
		if len(line) == 0 {
			continue
		}
		if len(line) >= 1 && line[0] == '#' {
			continue
		}
		if len(line) >= 2 && line[0] == '/' && line[1] == '/' {
			continue
		}
		kv := strings.Split(line, ": ")
		mode := kv[0]
		names := strings.Split(kv[1], ", ")
		for i, name := range names {
			names[i] = strings.TrimSpace(name)
		}
		keys := input.NamesToKeys(names)
		if !input.IsKeysValid(keys) {
			fmt.Printf("mapping keys are duplicated: %v\n", names)
			continue
		}
		switch mode {
		case "Drum", "drum":
			for _, prop := range props {
				if strings.Contains(strings.ToLower(prop.Name), "drum") {
					prop.KeySettings[4] = keys
					break
				}
			}
		default:
			subMode, err := strconv.Atoi(mode)
			if err != nil {
				fmt.Printf("error at loading key settings %s: %v", line, err)
				continue
			}
			for _, prop := range props {
				if strings.Contains(strings.ToLower(prop.Name), "piano") {
					prop.KeySettings[subMode] = keys
					break
				}
			}
		}
	}
}
