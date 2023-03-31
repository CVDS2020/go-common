package io

import (
	"os"
	"testing"
	"time"
)

func TestOpWriter(t *testing.T) {
	opWriter := OpWriter{}
	opWriter.AppendByte(1).
		AppendUint32(2).
		AppendUint16(3).
		AppendUint64(89).
		AppendBytes([]byte{0, 0, 0, 1, 67}).
		AppendUint16(3).
		AppendUint64(4).
		AppendString("hello world")
	OpWriterAppendEmbedded(&opWriter, time.Now())
	opWriter.WriteTo(os.Stdout)
}
