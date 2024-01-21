package install

import "time"

func loadingAnimation() string {
	animations := []string{"⌛", "⌛⌛", "⌛⌛⌛", "⌛⌛⌛⌛"}
	return animations[time.Now().Second()%4]
}
