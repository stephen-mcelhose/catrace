package stream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
)

// IngestJSONL streams JSON values from r (typically JSONL) into observe.
// Each non-empty line is decoded with encoding/json.Decoder; malformed lines
// are skipped so dirty mid-file dumps (as in some public samples) do not abort.
// Returns the number of observed non-empty states and skipped lines.
func IngestJSONL(r io.Reader, observe func(state string)) (observed, skipped int, err error) {
	br := bufio.NewReaderSize(r, 1<<20)
	for {
		line, readErr := br.ReadBytes('\n')
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) > 0 {
			var ev Event
			dec := json.NewDecoder(bytes.NewReader(trimmed))
			if decErr := dec.Decode(&ev); decErr != nil {
				skipped++
			} else if state := ev.CanonicalState(); state != "" {
				observe(state)
				observed++
			}
		}
		if readErr == io.EOF {
			return observed, skipped, nil
		}
		if readErr != nil {
			return observed, skipped, readErr
		}
	}
}
