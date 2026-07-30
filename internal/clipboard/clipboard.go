package clipboard

import "github.com/zyedidia/clipper"

var clip clipper.Clipboard

// Initialize clipboard
func Initialize() error {
	var err error
	clip, err = clipper.GetClipboard(clipper.Clipboards...)
	return err
}

// Returns contents of clipboard as string
func ReadAll() (string, error) {
	contents, err := clip.ReadAll(clipper.RegClipboard)
	return string(contents), err
}

// Writes to system clipboard
func WriteAll(data string) error {
	err := clip.WriteAll(clipper.RegClipboard, []byte(data))
	return err
}
