package render

import (
	"fmt"
	"io"
	"strings"
)

func AsTree(links []string, w io.Writer) {
	var prefix []string
	for _, link := range links {
		split := strings.Split(link, "/")[1:]
		depth := 0
		for i, path := range split {
			if i < len(prefix) && prefix[i] == path {
				depth++
			}
		}
		prefix = prefix[0:depth]

		for i := depth; i < len(split); i++ {
			path := split[i]
			for i := 0; i < 2*depth; i++ {
				fmt.Fprint(w, " ")
			}
			fmt.Fprintf(w, "- /%s\n", path)
			depth++
			prefix = append(prefix, path)
		}
	}
}
