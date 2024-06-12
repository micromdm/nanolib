package test

import "testing"

// TestLastLogKeyValueMatches verifies that a matching key-value pair is
// present in the last logged line of logger.
func TestLastLogKeyValueMatches(t *testing.T, logger *Logger, key, value string) {
	if logger == nil {
		panic("nil logger")
	}
	if key == "" {
		panic("empty key")
	}
	t.Run("key="+key, func(t *testing.T) {
		_, logValue, found := logger.LastKey(key)
		if !found {
			t.Fatal("key not found")
		}
		strLogValue, ok := logValue.(string)
		if !ok {
			t.Fatal("value not a string")
		}
		if have, want := strLogValue, value; have != want {
			t.Errorf("have = %q, want = %q", have, want)
		}
	})
}
