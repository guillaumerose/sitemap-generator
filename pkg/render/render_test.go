package render

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	links := []string{
		"/",
		"/en",
		"/en/about",
		"/en/about/company",
		"/en/resources",
		"/en/search",
		"/en/technologies",
		"/en/technologies/cloud-computing/cloud-suite",
		"/en/technologies/cloud-computing/openshift",
		"/en/topics/edge-computing/approach",
		"/en/websites-and-apps",
	}
	var buf bytes.Buffer
	AsTree(links, &buf)
	assert.Equal(t, `- /
- /en
  - /about
    - /company
  - /resources
  - /search
  - /technologies
    - /cloud-computing
      - /cloud-suite
      - /openshift
  - /topics
    - /edge-computing
      - /approach
  - /websites-and-apps
`, buf.String())
}
