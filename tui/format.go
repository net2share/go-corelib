package tui

// KV formats a key-value pair with theme colors.
func KV(key, value string) string {
	return KeyStyle.Render(key) + ValueStyle.Render(value)
}

// Header formats a header/section title.
func Header(text string) string {
	return HeaderStyle.Render(text)
}

// Code formats code/technical values like keys, passwords.
func Code(text string) string {
	return CodeStyle.Render(text)
}

// Value formats a standalone value with emphasis.
func Value(text string) string {
	return ValueStyle.Render(text)
}

// Muted formats text in a subtle, de-emphasized style.
func Muted(text string) string {
	return MutedStyle.Render(text)
}
