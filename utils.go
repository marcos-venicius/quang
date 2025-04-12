package quang

func isDigit[T byte | rune](c T) bool {
	return c >= '0' && c <= '9'
}

func isSymbol[T byte | rune](c T) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
