package pkg

import (
	"bytes"
	"io"
	"log"
)

const (
	bitsPerByte   = 8
	bytesPerint16 = 2
	bytesPerint32 = 4
	bytesPerint64 = 8
)

// streamReader allows you to read data from a byte array in various formats
type streamReader struct {
	data     []byte
	position uint64
}

// CreateStreamReader creates an instance of the stream reader
func CreateStreamReader(source []byte) *streamReader {
	result := &streamReader{
		data:     source,
		position: 0,
	}

	return result
}

// ReadByte reads a byte from the stream
func (v *streamReader) ReadByte() (byte, error) {
	if v.position >= v.Size() {
		return 0, io.EOF
	}

	result := v.data[v.position]
	v.position++

	return result, nil
}

// ReadInt16 returns a int16 word from the stream
func (v *streamReader) ReadInt16() (int16, error) {
	b, err := v.ReadUInt16()
	return int16(b), err
}

// ReadUInt16 returns a uint16 word from the stream
func (v *streamReader) ReadUInt16() (uint16, error) {
	b, err := v.ReadBytes(bytesPerint16)
	if err != nil {
		return 0, err
	}

	return uint16(b[0]) | uint16(b[1])<<8, err
}

// ReadInt32 returns an int32 dword from the stream
func (v *streamReader) ReadInt32() (int32, error) {
	b, err := v.ReadUInt32()
	return int32(b), err
}

// ReadUInt32 returns a uint32 dword from the stream
// nolint
func (v *streamReader) ReadUInt32() (uint32, error) {
	b, err := v.ReadBytes(bytesPerint32)
	if err != nil {
		return 0, err
	}

	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24, err
}

// ReadInt64 returns a uint64 qword from the stream
func (v *streamReader) ReadInt64() (int64, error) {
	b, err := v.ReadUInt64()
	return int64(b), err
}

// ReadUInt64 returns a uint64 qword from the stream
// nolint
func (v *streamReader) ReadUInt64() (uint64, error) {
	b, err := v.ReadBytes(bytesPerint64)
	if err != nil {
		return 0, err
	}

	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56, err
}

// Position returns the current stream position
func (v *streamReader) Position() uint64 {
	return v.position
}

// SetPosition sets the stream position with the given position
func (v *streamReader) SetPosition(newPosition uint64) {
	v.position = newPosition
}

// Size returns the total size of the stream in bytes
func (v *streamReader) Size() uint64 {
	return uint64(len(v.data))
}

// ReadBytes reads multiple bytes
func (v *streamReader) ReadBytes(count int) ([]byte, error) {
	if count <= 0 {
		return nil, nil
	}

	size := v.Size()
	if v.position >= size || v.position+uint64(count) > size {
		return nil, io.EOF
	}

	result := v.data[v.position : v.position+uint64(count)]
	v.position += uint64(count)

	return result, nil
}

// SkipBytes moves the stream position forward by the given amount
func (v *streamReader) SkipBytes(count int) {
	v.position += uint64(count)
}

// Read implements io.Reader
func (v *streamReader) Read(p []byte) (n int, err error) {
	streamLength := v.Size()

	for i := 0; ; i++ {
		if v.Position() >= streamLength {
			return i, io.EOF
		}

		if i >= len(p) {
			return i, nil
		}

		p[i], err = v.ReadByte()
		if err != nil {
			return i, err
		}
	}
}

// EOF returns if the stream position is reached to the end of the data, or not
func (v *streamReader) EOF() bool {
	return v.position >= uint64(len(v.data))
}

// streamWriter allows you to create a byte array by streaming in writes of various sizes
type streamWriter struct {
	data      *bytes.Buffer
	bitOffset int
	bitCache  byte
}

// CreateStreamWriter creates a new streamWriter instance
func CreateStreamWriter() *streamWriter {
	result := &streamWriter{
		data: new(bytes.Buffer),
	}

	return result
}

// GetBytes returns the the byte slice of the underlying data
func (v *streamWriter) GetBytes() []byte {
	return v.data.Bytes()
}

// PushBytes writes a bytes to the stream
func (v *streamWriter) PushBytes(b ...byte) {
	for _, i := range b {
		v.data.WriteByte(i)
	}
}

// PushBit pushes single bit into stream
// WARNING: if you'll use PushBit, offset'll be less than 8, and if you'll
// use another Push... method, bits'll not be pushed
func (v *streamWriter) PushBit(b bool) {
	if b {
		v.bitCache |= 1 << v.bitOffset
	}
	v.bitOffset++

	if v.bitOffset != bitsPerByte {
		return
	}

	v.PushBytes(v.bitCache)
	v.bitCache = 0
	v.bitOffset = 0
}

// PushBits pushes bits (with max range 8)
func (v *streamWriter) PushBits(b byte, bits int) {
	if bits > bitsPerByte {
		log.Print("input bits number must be less (or equal) than 8")
	}

	val := b

	for i := 0; i < bits; i++ {
		v.PushBit(val&1 == 1)
		val >>= 1
	}
}

// PushBits16 pushes bits (with max range 16)
func (v *streamWriter) PushBits16(b uint16, bits int) {
	if bits > bitsPerByte*bytesPerint16 {
		log.Print("input bits number must be less (or equal) than 16")
	}

	val := b

	for i := 0; i < bits; i++ {
		v.PushBit(val&1 == 1)
		val >>= 1
	}
}

// PushBits32 pushes bits (with max range 32)
func (v *streamWriter) PushBits32(b uint32, bits int) {
	if bits > bitsPerByte*bytesPerint32 {
		log.Print("input bits number must be less (or equal) than 32")
	}

	val := b

	for i := 0; i < bits; i++ {
		v.PushBit(val&1 == 1)
		val >>= 1
	}
}

// PushInt16 writes a int16 word to the stream
func (v *streamWriter) PushInt16(val int16) {
	v.PushUint16(uint16(val))
}

// PushUint16 writes an uint16 word to the stream
// nolint
func (v *streamWriter) PushUint16(val uint16) {
	v.data.WriteByte(byte(val))
	v.data.WriteByte(byte(val >> 8))
}

// PushInt32 writes a int32 dword to the stream
func (v *streamWriter) PushInt32(val int32) {
	v.PushUint32(uint32(val))
}

// PushUint32 writes a uint32 dword to the stream
// nolint
func (v *streamWriter) PushUint32(val uint32) {
	v.data.WriteByte(byte(val))
	v.data.WriteByte(byte(val >> 8))
	v.data.WriteByte(byte(val >> 16))
	v.data.WriteByte(byte(val >> 24))
}

// PushInt64 writes a uint64 qword to the stream
func (v *streamWriter) PushInt64(val int64) {
	v.PushUint64(uint64(val))
}

// PushUint64 writes a uint64 qword to the stream
// nolint
func (v *streamWriter) PushUint64(val uint64) {
	v.data.WriteByte(byte(val))
	v.data.WriteByte(byte(val >> 8))
	v.data.WriteByte(byte(val >> 16))
	v.data.WriteByte(byte(val >> 24))
	v.data.WriteByte(byte(val >> 32))
	v.data.WriteByte(byte(val >> 40))
	v.data.WriteByte(byte(val >> 48))
	v.data.WriteByte(byte(val >> 56))
}

const (
	maxBits = 16
)

// BitStream is a utility class for reading groups of bits from a stream
type BitStream struct {
	data         []byte
	dataPosition int
	current      int
	bitCount     int
}

// CreateBitStream creates a new BitStream
func CreateBitStream(newData []byte) *BitStream {
	result := &BitStream{
		data:         newData,
		dataPosition: 0,
		current:      0,
		bitCount:     0,
	}

	return result
}

// ReadBits reads the specified number of bits and returns the value
func (v *BitStream) ReadBits(bitCount int) int {
	if bitCount > maxBits {
		log.Panic("Maximum BitCount is 16")
	}

	if !v.EnsureBits(bitCount) {
		return -1
	}

	// nolint:gomnd // byte expresion
	result := v.current & (0xffff >> uint(maxBits-bitCount))
	v.WasteBits(bitCount)

	return result
}

// PeekByte returns the current byte without adjusting the position
func (v *BitStream) PeekByte() int {
	if !v.EnsureBits(bitsPerByte) {
		return -1
	}

	// nolint:gomnd // byte
	return v.current & 0xff
}

// EnsureBits ensures that the specified number of bits are available
func (v *BitStream) EnsureBits(bitCount int) bool {
	if bitCount <= v.bitCount {
		return true
	}

	if v.dataPosition >= len(v.data) {
		return false
	}

	nextValue := v.data[v.dataPosition]
	v.dataPosition++
	v.current |= int(nextValue) << uint(v.bitCount)
	v.bitCount += 8

	return true
}

// WasteBits dry-reads the specified number of bits
func (v *BitStream) WasteBits(bitCount int) {
	// noinspection GoRedundantConversion
	v.current >>= uint(bitCount)
	v.bitCount -= bitCount
}
