package core

// Hello returns a friendly greeting.
func Hello(name string) string {
	if name == "" {
		return "Hello, World!"
	}
	return "Hello, " + name + "!"
}
