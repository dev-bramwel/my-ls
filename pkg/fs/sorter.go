package fs

func SortFiles(files []FileInfo, timeSort bool, reverse bool) []FileInfo {
	if len(files) <= 1 {
		return files
	}

	n := len(files)

	for i := 0; i < n-1; i++ {
		extreme := i

		for j := i + 1; j < n; j++ {
			if timeSort {
				if files[j].ModTime.Equal(files[extreme].ModTime) {
					nameJ := ToLower(files[j].Name)
					nameExt := ToLower(files[extreme].Name)

					if !reverse {
						if nameJ < nameExt || (nameJ == nameExt && files[j].Name < files[extreme].Name) {
							extreme = j
						}
					} else {
						if nameJ > nameExt || (nameJ == nameExt && files[j].Name > files[extreme].Name) {
							extreme = j
						}
					}
				} else if !reverse {
					if files[j].ModTime.After(files[extreme].ModTime) {
						extreme = j
					}
				} else {
					if files[j].ModTime.Before(files[extreme].ModTime) {
						extreme = j
					}
				}
			} else {
				nameJ := ToLower(files[j].Name)
				nameExt := ToLower(files[extreme].Name)

				if !reverse {
					if nameJ < nameExt || (nameJ == nameExt && files[j].Name < files[extreme].Name) {
						extreme = j
					}
				} else {
					if nameJ > nameExt || (nameJ == nameExt && files[j].Name > files[extreme].Name) {
						extreme = j
					}
				}
			}
		}

		if extreme != i {
			files[i], files[extreme] = files[extreme], files[i]
		}
	}

	return files
}

func ToLower(s string) string {
	if len(s) > 1 && s[0] == '.' && s != ".." {
		s = s[1:]
	}
	b := []byte(s)
	for i := 0; i < len(b); i++ {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 32
		}
	}
	return string(b)
}