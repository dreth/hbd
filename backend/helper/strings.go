package helper

// joinStrings joins the strings with a specified separator.
func JoinStrings(elements []string, separator string) string {
	var result string
	for i, element := range elements {
		if i > 0 {
			result += separator
		}
		result += element
	}
	return result
}
