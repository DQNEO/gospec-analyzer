package wordprocessor

import "strings"

func Singulify(word string) string {
	var r string
	switch {
	case strings.HasSuffix(word, "ies"):
		// "entries" => "entry"
		r = strings.TrimSuffix(word, "ies") + "y"
	case strings.HasSuffix(word, "es"):
		// @TODO "resumes" => "resume"
		// "classes" => "class"
		r = strings.TrimSuffix(word, "es")
	default:
		r = strings.TrimSuffix(word, "s")
	}
	return r
}
