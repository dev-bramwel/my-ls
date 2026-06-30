package fs

// SortFiles sorts a slice of FileInfo structs based on time and reverse flags.
// Parameters:
//   - files: slice of FileInfo to sort (modified in place)
//   - timeSort: when true, sort by modification time; when false, sort alphabetically by name
//   - reverse: when true, reverse the sort order (ls default: time desc, name asc)
//
// Returns:
//   - []FileInfo: the sorted slice (same underlying data, reordered)
//
// Scope: Performs in-place sorting using selection sort algorithm.
// Flag interpretation:
//   - Default (no flags): alphabetical ascending
//   - -t only: time descending (newest first)
//   - -r only: alphabetical descending
//   - -t -r or -rt: time ascending (oldest first)
//   - -r -t: same as -t -r (time ascending)
func SortFiles(files []FileInfo, timeSort bool, reverse bool) []FileInfo {
	if len(files) <= 1 {
		return files
	}

	// Use selection sort - O(n²) but no external dependencies needed
	n := len(files)

	for i := 0; i < n-1; i++ {
		extreme := i

		for j := i + 1; j < n; j++ {
			_ = false // suppress unused variable warning

			if timeSort {
				// Sort by modification time
				// Default ls -t behavior: newest first (descending)
				if !reverse {
					// Descending: find file with LATER time
					if files[j].ModTime.After(files[extreme].ModTime) {
						extreme = j
					}
				} else {
					// Ascending: find file with EARLIER time
					if files[j].ModTime.Before(files[extreme].ModTime) {
						extreme = j
					}
				}
			} else {
				// Sort alphabetically by name
				// Case-insensitive, stripping leading dots to match standard ls -a positioning
				nameJ := toLower(files[j].Name)
				nameExt := toLower(files[extreme].Name)

				if !reverse {
					// Ascending: find file with "smaller" name
					if nameJ < nameExt || (nameJ == nameExt && files[j].Name < files[extreme].Name) {
						extreme = j
					}
				} else {
					// Descending: find file with "larger" name
					if nameJ > nameExt || (nameJ == nameExt && files[j].Name > files[extreme].Name) {
						extreme = j
					}
				}
			}
		}

		// Swap the extreme element with position i
		if extreme != i {
			files[i], files[extreme] = files[extreme], files[i]
		}
	}

	return files
}

// toLower converts an ASCII string to lowercase while stripping a leading dot for correct ls -a sorting.
func toLower(s string) string {
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