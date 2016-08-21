// Package jwriter contains a JSON writer.
package jwriter

import (
	"io"
	"strconv"
	"unicode/utf8"
	"sync"

	bbp "github.com/valyala/bytebufferpool"
)

var (
	wPool = &sync.Pool{
		New: func() interface{} {
			return &Writer{nil, bbp.Get()}
		},
	}
)

// Writer is a JSON writer.
type Writer struct {
	Error  error
	Buffer *bbp.ByteBuffer
}

func New() *Writer {
	return wPool.Get().(*Writer)
}

func Free(w *Writer) {
	w.Buffer.Reset()
	wPool.Put(w)
}

// Size returns the size of the data that was written out.
func (w *Writer) Size() int {
	return w.Buffer.Len()
}

// DumpTo outputs the data to given io.Writer, resetting the buffer.
func (w *Writer) DumpTo(out io.Writer) (written int64, err error) {
	return w.Buffer.WriteTo(out)
}

// BuildBytes returns writer data as a single byte slice.
func (w *Writer) BuildBytes() ([]byte, error) {
	if w.Error != nil {
		return nil, w.Error
	}

	return w.Buffer.Bytes(), nil
}

// RawByte appends raw binary data to the buffer.
func (w *Writer) RawByte(c byte) {
	w.Buffer.WriteByte(c)
}

// RawByte appends raw binary data to the buffer.
func (w *Writer) RawString(s string) {
	w.Buffer.WriteString(s)
}

// RawByte appends raw binary data to the buffer or sets the error if it is given. Useful for
// calling with results of MarshalJSON-like functions.
func (w *Writer) Raw(data []byte, err error) {
	switch {
	case w.Error != nil:
		return
	case err != nil:
		w.Error = err
	case len(data) > 0:
		w.Buffer.Write(data)
	default:
		w.RawString("null")
	}
}

func (w *Writer) Uint8(n uint8) {
	w.Buffer.B = strconv.AppendUint(w.Buffer.B, uint64(n), 10)
}

func (w *Writer) Uint16(n uint16) {
	w.Buffer.B = strconv.AppendUint(w.Buffer.B, uint64(n), 10)
}

func (w *Writer) Uint32(n uint32) {
	w.Buffer.B = strconv.AppendUint(w.Buffer.B, uint64(n), 10)
}

func (w *Writer) Uint(n uint) {
	w.Buffer.B = strconv.AppendUint(w.Buffer.B, uint64(n), 10)
}

func (w *Writer) Uint64(n uint64) {
	w.Buffer.B = strconv.AppendUint(w.Buffer.B, n, 10)
}

func (w *Writer) Int8(n int8) {
	w.Buffer.B = strconv.AppendInt(w.Buffer.B, int64(n), 10)
}

func (w *Writer) Int16(n int16) {
	w.Buffer.B = strconv.AppendInt(w.Buffer.B, int64(n), 10)
}

func (w *Writer) Int32(n int32) {
	w.Buffer.B = strconv.AppendInt(w.Buffer.B, int64(n), 10)
}

func (w *Writer) Int(n int) {
	w.Buffer.B = strconv.AppendInt(w.Buffer.B, int64(n), 10)
}

func (w *Writer) Int64(n int64) {
	w.Buffer.B = strconv.AppendInt(w.Buffer.B, n, 10)
}

func (w *Writer) Uint8Str(n uint8) {
	w.Buffer.WriteByte('"')
	w.Uint8(n)
	w.Buffer.WriteByte('"')
}

func (w *Writer) Uint16Str(n uint16) {
	w.Buffer.WriteByte('"')
	w.Uint16(n)
	w.Buffer.WriteByte('"')
}

func (w *Writer) Uint32Str(n uint32) {
	w.Buffer.WriteByte('"')
	w.Uint32(n)
	w.Buffer.WriteByte('"')
}

func (w *Writer) UintStr(n uint) {
	w.Buffer.WriteByte('"')
	w.Uint(n)
	w.Buffer.WriteByte('"')
}

func (w *Writer) Uint64Str(n uint64) {
	w.Buffer.WriteByte('"')
	w.Uint64(n)
	w.Buffer.WriteByte('"')
}

func (w *Writer) Int8Str(n int8) {
	w.Buffer.WriteByte('"')
	w.Int8(n)
	w.Buffer.WriteByte('"')
}

func (w *Writer) Int16Str(n int16) {
	w.Buffer.WriteByte('"')
	w.Int16(n)
	w.Buffer.WriteByte('"')
}

func (w *Writer) Int32Str(n int32) {
	w.Buffer.WriteByte('"')
	w.Int32(n)
	w.Buffer.WriteByte('"')
}

func (w *Writer) IntStr(n int) {
	w.Buffer.WriteByte('"')
	w.Int(n)
	w.Buffer.WriteByte('"')
}

func (w *Writer) Int64Str(n int64) {
	w.Buffer.WriteByte('"')
	w.Int64(n)
	w.Buffer.WriteByte('"')
}

func (w *Writer) Float32(n float32) {
	w.Buffer.B = strconv.AppendFloat(w.Buffer.B, float64(n), 'g', -1, 32)
}

func (w *Writer) Float64(n float64) {
	w.Buffer.B = strconv.AppendFloat(w.Buffer.B, n, 'g', -1, 64)
}

func (w *Writer) Bool(v bool) {
	if v {
		w.Buffer.WriteString("true")
	} else {
		w.Buffer.WriteString("false")
	}
}

const chars = "0123456789abcdef"

func isNotEscapedSingleChar(c byte) bool {
	// Note: might make sense to use a table if there are more chars to escape. With 4 chars
	// it benchmarks the same.
	return c != '<' && c != '\\' && c != '"' && c != '>' && c >= 0x20 && c < utf8.RuneSelf
}

func (w *Writer) String(s string) {
	w.Buffer.WriteByte('"')

	// Portions of the string that contain no escapes are appended as
	// byte slices.

	p := 0 // last non-escape symbol

	for i := 0; i < len(s); {
		c := s[i]

		if isNotEscapedSingleChar(c) {
			// single-width character, no escaping is required
			i++
			continue
		} else if c < utf8.RuneSelf {
			// single-with character, need to escape
			w.Buffer.WriteString(s[p:i])
			switch c {
			case '\t':
				w.Buffer.WriteString(`\t`)
			case '\r':
				w.Buffer.WriteString(`\r`)
			case '\n':
				w.Buffer.WriteString(`\n`)
			case '\\':
				w.Buffer.WriteString(`\\`)
			case '"':
				w.Buffer.WriteString(`\"`)
			default:
				w.Buffer.WriteString(`\u00`)
				w.Buffer.WriteByte(chars[c>>4])
				w.Buffer.WriteByte(chars[c&0xf])
			}

			i++
			p = i
			continue
		}

		// broken utf
		runeValue, runeWidth := utf8.DecodeRuneInString(s[i:])
		if runeValue == utf8.RuneError && runeWidth == 1 {
			w.Buffer.WriteString(s[p:i])
			w.Buffer.WriteString(`\ufffd`)
			i++
			p = i
			continue
		}

		// jsonp stuff - tab separator and line separator
		if runeValue == '\u2028' || runeValue == '\u2029' {
			w.Buffer.WriteString(s[p:i])
			w.Buffer.WriteString(`\u202`)
			w.Buffer.WriteByte(chars[runeValue&0xf])
			i += runeWidth
			p = i
			continue
		}
		i += runeWidth
	}
	w.Buffer.WriteString(s[p:])
	w.Buffer.WriteByte('"')
}
