/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/distroy/ldgo/v3/ldlog/internal/handler"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// For JSON-escaping; see jsonEncoder.safeAddString below.
const (
	_hex = "0123456789abcdef"
)

type loggerEncoderPool struct {
	_pool sync.Pool
}

func (pool *loggerEncoderPool) get() *loggerEncoder {
	return pool._pool.Get().(*loggerEncoder)
}

func (pool *loggerEncoderPool) put(enc *loggerEncoder) {
	if enc.reflectBuf != nil {
		enc.reflectBuf.Free()
	}
	enc.EncoderConfig = nil
	enc.buf = nil
	enc.sequence = ""
	enc.spaced = false
	enc.openNamespaces = 0
	enc.reflectBuf = nil
	enc.reflectEnc = nil
	pool._pool.Put(enc)
}

var (
	bufferpool  = buffer.NewPool()
	encoderPool = &loggerEncoderPool{
		_pool: sync.Pool{
			New: func() interface{} {
				return &loggerEncoder{}
			},
		},
	}
)

type loggerEncoder struct {
	*zapcore.EncoderConfig
	buf            *buffer.Buffer
	spaced         bool // include spaces after colons and commas
	openNamespaces int
	// for encoding generic values by reflection
	reflectBuf *buffer.Buffer
	reflectEnc *json.Encoder

	sequence string
}

// NewLoggerEncoder
// log format: 2013-04-08 15:30:42.621|FATAL|0x7f865e4b9720|GdpProcessor.cpp(35)|CGdpPrcessor::Init|init_db_client_fail|id=95,db=gpp_db,type=1
func NewLoggerEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &loggerEncoder{
		EncoderConfig: &cfg,
		buf:           bufferpool.Get(),
		spaced:        false,
	}
}

func (enc *loggerEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	enc.addKey(key)
	return enc.AppendArray(arr)
}

func (enc *loggerEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

func (enc *loggerEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *loggerEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
}

func (enc *loggerEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *loggerEncoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
}

func (enc *loggerEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
}

func (enc *loggerEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
}

func (enc *loggerEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
}

func (enc *loggerEncoder) resetReflectBuf() {
	if enc.reflectBuf == nil {
		enc.reflectBuf = bufferpool.Get()
		enc.reflectEnc = json.NewEncoder(enc.reflectBuf)
	} else {
		enc.reflectBuf.Reset()
	}
}

var nullLiteralBytes = []byte("null")

func (enc *loggerEncoder) encodeReflected(obj interface{}) ([]byte, error) {
	if obj == nil {
		return nullLiteralBytes, nil
	}
	enc.resetReflectBuf()
	if err := enc.reflectEnc.Encode(obj); err != nil {
		return nil, err
	}
	enc.reflectBuf.TrimNewline()
	return enc.reflectBuf.Bytes(), nil
}

func (enc *loggerEncoder) AddReflected(key string, obj interface{}) error {
	val, err := enc.encodeReflected(obj)
	if err != nil {
		return err
	}
	enc.addKey(key)
	_, err = enc.buf.Write(val)
	return err
}

func (enc *loggerEncoder) OpenNamespace(key string) {
	enc.addKey(key)
	enc.buf.AppendByte('{')
	enc.openNamespaces++
}

func (enc *loggerEncoder) AddString(key, val string) {
	switch key {
	case handler.SequenceKey:
		enc.sequence = val
	default:
		enc.addKey(key)
		enc.AppendString(val)
	}
}

func (enc *loggerEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTime(val)
}

func (enc *loggerEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
}

func (enc *loggerEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('[')
	err := arr.MarshalLogArray(enc)
	enc.buf.AppendByte(']')
	return err
}

func (enc *loggerEncoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('{')
	err := obj.MarshalLogObject(enc)
	enc.buf.AppendByte('}')
	return err
}

func (enc *loggerEncoder) AppendBool(val bool) {
	enc.addElementSeparator()
	enc.buf.AppendBool(val)
}

func (enc *loggerEncoder) AppendByteString(val []byte) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddByteString(val)
	enc.buf.AppendByte('"')
}

func (enc *loggerEncoder) AppendComplex128(val complex128) {
	enc.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *loggerEncoder) AppendDuration(val time.Duration) {
	cur := enc.buf.Len()
	enc.EncodeDuration(val, enc)
	if cur == enc.buf.Len() {
		// User-supplied EncodeDuration is a no-op. Fall back to nanoseconds to keep
		// JSON valid.
		enc.AppendInt64(int64(val))
	}
}

func (enc *loggerEncoder) AppendInt64(val int64) {
	enc.addElementSeparator()
	enc.buf.AppendInt(val)
}

func (enc *loggerEncoder) AppendReflected(obj interface{}) error {
	val, err := enc.encodeReflected(obj)
	if err != nil {
		return err
	}
	enc.addElementSeparator()
	_, err = enc.buf.Write(val)
	return err
}

func (enc *loggerEncoder) AppendString(val string) {
	enc.addElementSeparator()
	enc.safeAddString(val)
}

func (enc *loggerEncoder) AppendTime(val time.Time) {
	cur := enc.buf.Len()
	enc.EncodeTime(val, enc)
	if cur == enc.buf.Len() {
		// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
		// output JSON valid.
		enc.AppendInt64(val.UnixNano())
	}
}

func (enc *loggerEncoder) AppendUint64(val uint64) {
	enc.addElementSeparator()
	enc.buf.AppendUint(val)
}

func (enc *loggerEncoder) AddComplex64(k string, v complex64) { enc.AddComplex128(k, complex128(v)) }
func (enc *loggerEncoder) AddFloat32(k string, v float32)     { enc.AddFloat64(k, float64(v)) }
func (enc *loggerEncoder) AddInt(k string, v int)             { enc.AddInt64(k, int64(v)) }
func (enc *loggerEncoder) AddInt32(k string, v int32)         { enc.AddInt64(k, int64(v)) }
func (enc *loggerEncoder) AddInt16(k string, v int16)         { enc.AddInt64(k, int64(v)) }
func (enc *loggerEncoder) AddInt8(k string, v int8)           { enc.AddInt64(k, int64(v)) }
func (enc *loggerEncoder) AddUint(k string, v uint)           { enc.AddUint64(k, uint64(v)) }
func (enc *loggerEncoder) AddUint32(k string, v uint32)       { enc.AddUint64(k, uint64(v)) }
func (enc *loggerEncoder) AddUint16(k string, v uint16)       { enc.AddUint64(k, uint64(v)) }
func (enc *loggerEncoder) AddUint8(k string, v uint8)         { enc.AddUint64(k, uint64(v)) }
func (enc *loggerEncoder) AddUintptr(k string, v uintptr)     { enc.AddUint64(k, uint64(v)) }
func (enc *loggerEncoder) AppendComplex64(v complex64)        { enc.AppendComplex128(complex128(v)) }
func (enc *loggerEncoder) AppendFloat64(v float64)            { enc.appendFloat(v, 64) }
func (enc *loggerEncoder) AppendFloat32(v float32)            { enc.appendFloat(float64(v), 32) }
func (enc *loggerEncoder) AppendInt(v int)                    { enc.AppendInt64(int64(v)) }
func (enc *loggerEncoder) AppendInt32(v int32)                { enc.AppendInt64(int64(v)) }
func (enc *loggerEncoder) AppendInt16(v int16)                { enc.AppendInt64(int64(v)) }
func (enc *loggerEncoder) AppendInt8(v int8)                  { enc.AppendInt64(int64(v)) }
func (enc *loggerEncoder) AppendUint(v uint)                  { enc.AppendUint64(uint64(v)) }
func (enc *loggerEncoder) AppendUint32(v uint32)              { enc.AppendUint64(uint64(v)) }
func (enc *loggerEncoder) AppendUint16(v uint16)              { enc.AppendUint64(uint64(v)) }
func (enc *loggerEncoder) AppendUint8(v uint8)                { enc.AppendUint64(uint64(v)) }
func (enc *loggerEncoder) AppendUintptr(v uintptr)            { enc.AppendUint64(uint64(v)) }

func (enc *loggerEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	clone.buf.Write(enc.buf.Bytes()) // nolint
	return clone
}

func (enc *loggerEncoder) clone() *loggerEncoder {
	clone := encoderPool.get()
	clone.sequence = enc.sequence
	clone.EncoderConfig = enc.EncoderConfig
	clone.spaced = enc.spaced
	clone.openNamespaces = enc.openNamespaces
	clone.buf = bufferpool.Get()
	return clone
}

func (enc *loggerEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	clone := enc.clone()
	addTimestamp(clone, &ent)
	addLogLevel(clone, &ent)
	addRequestID(clone)
	addLoggerName(clone, &ent)

	addCaller(clone, &ent)
	addMsg(clone, &ent)
	addBuffer(enc, clone)
	addFields(clone, fields)
	clone.closeOpenNamespaces()
	addStacktrace(clone, &ent)
	addLineEnding(clone)
	ret := clone.buf
	encoderPool.put(clone)
	return ret, nil
}

// nolint
func (enc *loggerEncoder) truncate() {
	enc.buf.Reset()
}

func (enc *loggerEncoder) closeOpenNamespaces() {
	for i := 0; i < enc.openNamespaces; i++ {
		enc.buf.AppendByte('}')
	}
}

func (enc *loggerEncoder) addKey(key string) {
	enc.addElementSeparator()
	enc.safeAddString(key)
	enc.buf.AppendByte('=')
	if enc.spaced {
		enc.buf.AppendByte(' ')
	}
}

func (enc *loggerEncoder) addElementSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}
	switch enc.buf.Bytes()[last] {
	case '|', '=', '{', '[', ':', ',', ' ':
		return
	default:
		enc.buf.AppendByte(',')
		if enc.spaced {
			enc.buf.AppendByte(' ')
		}
	}
}

func (enc *loggerEncoder) addLogSegmentSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}
	if enc.buf.Bytes()[last] == '|' {
		return
	}

	enc.buf.AppendByte('|')
}

func (enc *loggerEncoder) appendFloat(val float64, bitSize int) {
	enc.addElementSeparator()
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *loggerEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *loggerEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size]) // nolint
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *loggerEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *loggerEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func addTimestamp(enc *loggerEncoder, ent *zapcore.Entry) {
	if enc == nil || ent == nil {
		return
	}
	if enc.TimeKey != "" {
		// s := enc.E
		enc.EncodeTime(ent.Time, enc)
		// enc.buf.AppendString(ent.Time.Format("2006-01-02 15:04:05.000"))
	}
}
func addLogLevel(enc *loggerEncoder, ent *zapcore.Entry) {
	if enc == nil || ent == nil {
		return
	}
	if enc.LevelKey != "" {
		enc.addLogSegmentSeparator()
		enc.buf.AppendString(strings.ToUpper(ent.Level.String()))
	}
}

func addLoggerName(enc *loggerEncoder, ent *zapcore.Entry) {
	if ent.LoggerName != "" && enc.NameKey != "" {
		enc.addKey(enc.NameKey)
		cur := enc.buf.Len()
		nameEncoder := enc.EncodeName

		// if no name encoder provided, fall back to FullNameEncoder for backwards
		// compatibility
		if nameEncoder == nil {
			nameEncoder = zapcore.FullNameEncoder
		}

		nameEncoder(ent.LoggerName, enc)
		if cur == enc.buf.Len() {
			// User-supplied EncodeName was a no-op. Fall back to strings to
			// keep output JSON valid.
			enc.AppendString(ent.LoggerName)
		}
	}
}

func addMsg(enc *loggerEncoder, ent *zapcore.Entry) {
	if enc.MessageKey != "" {
		enc.addLogSegmentSeparator()
		enc.AppendString(ent.Message)
	}
}

func addCaller(enc *loggerEncoder, ent *zapcore.Entry) {
	if ent.Caller.Defined && enc.CallerKey != "" {
		enc.addLogSegmentSeparator()
		cur := enc.buf.Len()
		enc.EncodeCaller(ent.Caller, enc)
		if cur == enc.buf.Len() {
			// User-supplied EncodeCaller was a no-op. Fall back to strings to
			// keep output JSON valid.
			enc.AppendString(ent.Caller.String())
		}
		return
	}
	enc.buf.AppendString("|-")

}

func addBuffer(enc, final *loggerEncoder) {
	if enc.buf.Len() > 0 {
		final.addElementSeparator()
		final.buf.Write(enc.buf.Bytes()) // nolint
	}
}

// func add
func addRequestID(enc *loggerEncoder) {
	var requestID = "-"
	if enc.sequence != "" {
		requestID = enc.sequence
	}
	enc.addLogSegmentSeparator()
	enc.buf.AppendString(requestID)
}

// func add
func addFields(enc *loggerEncoder, fields []zapcore.Field) {
	if len(fields) == 0 {
		return
	}
	enc.addLogSegmentSeparator()
	for i := range fields {
		fields[i].AddTo(enc)
	}
}
func addStacktrace(enc *loggerEncoder, ent *zapcore.Entry) {
	if ent.Stack != "" && enc.StacktraceKey != "" {
		enc.AddString(enc.StacktraceKey, ent.Stack)
	}
}
func addLineEnding(enc *loggerEncoder) {
	if enc.LineEnding != "" {
		enc.buf.AppendString(enc.LineEnding)
	} else {
		enc.buf.AppendString(zapcore.DefaultLineEnding)
	}
}
