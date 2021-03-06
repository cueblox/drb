package hosting

import (
	"fmt"
	"os"
)

func CreateFileWithContents(path string, contents string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating file: %s", err)
	}
	defer f.Close()

	_, err = f.WriteString(contents)
	return err
}
