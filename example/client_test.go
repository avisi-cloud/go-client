package example

import "testing"

func TestRunExample(t *testing.T) {
	// test/example is disabled due to missing
	// - ACLOUD_PAT
	// - ACLOUD_URL
	t.SkipNow()

	RunExample()
}
