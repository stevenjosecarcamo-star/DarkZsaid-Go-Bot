package bot

import "strings"

func cleanBrailleText(s string) string {
var b strings.Builder
lastNewline := false
blankCount := 0

for _, r := range s {
if r >= '\u2800' && r <= '\u28FF' {
continue
}

if r == '\n' {
if lastNewline {
blankCount++
if blankCount > 1 {
continue
}
} else {
blankCount = 0
}
lastNewline = true
} else if r != ' ' && r != '\t' && r != '\r' {
lastNewline = false
}

b.WriteRune(r)
}

return strings.TrimSpace(b.String())
}
