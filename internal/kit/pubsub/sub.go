package pubsub

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func DecodeMessageData(source string, dst any) error {
	data, err := base64.StdEncoding.DecodeString(source)
	if err != nil {
		return fmt.Errorf("cannot decode message data: %w", err)
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("cannot unmarshal JSON into struct: %w", err)
	}

	return nil
}
