package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// stringPrompt asks for a string value using the label
func stringPrompt(label, def string, required bool) string {
	text := fmt.Sprintf("[+] %s: ", label)
	if len(def) > 0 {
		text = fmt.Sprintf("[+] %s [%s]: ", label, def)
	}

	fmt.Fprint(os.Stderr, text)

	in := bufio.NewReader(os.Stdin)

	ans, _ := in.ReadString('\n')
	ans = strings.TrimSpace(ans)
	if ans == "" {
		ans = def
	}

	if ans != "" || !required {
		return ans
	}

	return stringPrompt(label, def, required)
}

// yesNoPrompt asks yes/no questions using the label.
func yesNoPrompt(label string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}
	text := fmt.Sprintf("[+] %s [%s]: ", label, choices)

	fmt.Fprint(os.Stderr, text)

	in := bufio.NewReader(os.Stdin)

	ans, _ := in.ReadString('\n')
	ans = strings.TrimSpace(ans)
	if ans == "" {
		return def
	}

	ans = strings.ToLower(ans)
	return (ans == "y" || ans == "yes")
}
