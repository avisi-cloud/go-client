package example

import "testing"

func TestRunAdminExample(t *testing.T) {
	// test/example is disabled due to missing
	// - ACLOUD_PAT
	// - ACLOUD_URL
	t.SkipNow()

	RunAdminExample()
}
