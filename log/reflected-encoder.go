package log

import (
	"encoding/json"
	"io"
)

// ReflectedEncoder serializes log fields that can't be serialized with Zap's
// JSON encoder. These have the ReflectType field type.
// Use EncoderConfig.NewReflectedEncoder to set this.
type ReflectedEncoder interface {
	// Encode encodes and writes to the underlying data stream.
	Encode(interface{}) error
}

func defaultReflectedEncoder(w io.Writer) ReflectedEncoder {
	enc := json.NewEncoder(w)
	// For consistency with our custom JSON encoder.
	enc.SetEscapeHTML(false)
	return enc
}
