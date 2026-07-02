package fs

// SortFiles leverages selection sort to arrange file slices dynamically based on flag variations
func SortFiles(files []FileInfo, timeSort bool, reverse bool) []FileInfo {
	if len(files) <= 1 {
		return files
	}

	n := len(files) // n tracks item scope metrics within selection sort loops

	for i := 0; i < n-1; i++ {
		extreme := i // extreme tracks index locations matching the target sorted position

		for j := i + 1; j < n; j++ {
			if timeSort {
				// Primary sorting logic evaluating modification time metrics (-t)
				if files[j].ModTime.Equal(files[extreme].ModTime) {
					// Fall back to alphabetical name comparisons if modification timestamps match exactly
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
					// Newest files first (descending timestamp metrics)
					if files[j].ModTime.After(files[extreme].ModTime) {
						extreme = j
					}
				} else {
					// Oldest files first (ascending timestamp metrics matching active -r option)
					if files[j].ModTime.Before(files[extreme].ModTime) {
						extreme = j
					}
				}
			} else {
				// Standard sorting logic matching regular alphabetical name sorting
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

// ToLower isolates and maps lowercase transformation vectors while ignoring leading dots on hidden targets
func ToLower(s string) string {
	// Strip the leading dot from hidden targets (like .config) to mimic native ls sequence listings
	if len(s) > 1 && s[0] == '.' && s != ".." {
		s = s[1:]
	}
	b := []byte(s) // b creates a mutable byte array tracking internal translation loops
	for i := 0; i < len(b); i++ {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 32 // Offset mapping constant shifting uppercase ASCII straight down into lowercase
		}
	}
	return string(b)
}