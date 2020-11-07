package duration_counter

import (
	"fmt"
	"io"
	"os"

	"github.com/tcolgate/mp3"
)

func CountDuration(path string) (int, error) {
	filePtr, err := os.Open(path)

	if err != nil {
		return 0, err
	}

	decoder := mp3.NewDecoder(filePtr)

	if err != nil {
		return 0, err
	}

	var f mp3.Frame
	duration := 0.0
	skipped := 0

	for {

		if err := decoder.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return 0, err
		}

		duration += f.Duration().Seconds()
	}

	return int(duration), nil
}
