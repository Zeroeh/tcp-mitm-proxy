package main

import (
	"math"
	"encoding/binary"
)

//Packet is a struct containing an underlying bytes buffer
type Packet struct {
	Index  int
	Length int
	ID     int
	Data   []byte
}

//Advance the buffer index while returning the amount increased
func (p *Packet) Advance(amount int) int {
	p.Index += amount
	return amount
}

//ReadString reads the expected string size (n) and read until n
func (p *Packet)ReadString() string {
	n := int(p.ReadUInt16()) //absolute
	if n == 0 {
		return ""
	}
	var str []byte
	str = p.Data[p.Index:p.Index+n]
	p.Advance(n)
	return string(str)
}

//WriteString writes int16 (len of string) and then the contents of s as bytes
func (p *Packet)WriteString(s string) {
	if s == "" {
		p.WriteUInt16(uint16(0))
		return
	}
	p.WriteUInt16(uint16(len(s)))
	for i := range s {
		p.WriteByte(s[i])
	}
}

func (p *Packet)ReadUTFString() string {
	n := int(p.ReadUInt32())
	if n == 0 {
		return ""
	}
	var str []byte
	str = p.Data[p.Index:p.Index+n]
	p.Advance(n)
	return string(str)
}

func (p *Packet)WriteUTFString(s string) {
	if s == "" {
		p.WriteUInt32(0)
		return
	}
	p.WriteUInt32(uint32(len(s)))
	for i := range s {
		p.WriteByte(s[i])
	}
}

func (p *Packet)ReadBool() bool {
	if p.ReadByte() == 1 {
		return true
	}
	return false //assume anything else is false
}

func (p *Packet)WriteBool(b bool) {
	if b == true {
		p.WriteByte(1)
	} else {
		p.WriteByte(0)
	}
}

//ReadFloat reads 4 bytes representing a float
func (p *Packet)ReadFloat() float32 {
	return math.Float32frombits(p.ReadUInt32())
}

//WriteFloat writes 4 bytes representing a float
func (p *Packet)WriteFloat(f float32) {
	binary.BigEndian.PutUint32(p.Data[p.Index:p.Index+p.Advance(4)], math.Float32bits(f))
}

func (p *Packet)ReadInt16() int16 {
	return int16(binary.BigEndian.Uint16(p.Data[p.Index:p.Index+p.Advance(2)]))
}

func (p *Packet)WriteInt16(i int16) {
	binary.BigEndian.PutUint16(p.Data[p.Index:p.Index+p.Advance(2)], uint16(i))
}

func (p *Packet)ReadUInt16() uint16 {
	return binary.BigEndian.Uint16(p.Data[p.Index:p.Index+p.Advance(2)])
}

func (p *Packet)WriteUInt16(i uint16) {
	binary.BigEndian.PutUint16(p.Data[p.Index:p.Index+p.Advance(2)], i)
}

func (p *Packet)ReadInt32() int32 {
	return int32(binary.BigEndian.Uint32(p.Data[p.Index:p.Index+p.Advance(4)]))
}

func (p *Packet)WriteInt32(i int32) {
	binary.BigEndian.PutUint32(p.Data[p.Index:p.Index+p.Advance(4)], uint32(i))
}

func (p *Packet)ReadUInt32() uint32 {
	return binary.BigEndian.Uint32(p.Data[p.Index:p.Index+p.Advance(4)])
}

func (p *Packet)WriteUInt32(i uint32) {
	binary.BigEndian.PutUint32(p.Data[p.Index:p.Index+p.Advance(4)], i)
}

func (p *Packet) WriteInt64(i int64) {
	binary.BigEndian.PutUint64(p.Data[p.Index:p.Index+8], uint64(i))
	p.Advance(8)
}

func (p *Packet) ReadInt64() int64 {
	b := int64(binary.BigEndian.Uint64(p.Data[p.Index : p.Index+8]))
	p.Advance(8)
	return b
}

func (p *Packet) WriteUInt64(i uint64) {
	binary.BigEndian.PutUint64(p.Data[p.Index:p.Index+8], i)
	p.Advance(8)
}

func (p *Packet) ReadUInt64() uint64 {
	b := binary.BigEndian.Uint64(p.Data[p.Index : p.Index+8])
	p.Advance(8)
	return b
}

//ReadByte reads and returns a singular byte
func (p *Packet)ReadByte() byte {
	return p.Data[p.Index:p.Index+p.Advance(1)][0]
}

//WriteByte writes a singular byte to the packet buffer
func (p *Packet)WriteByte(d byte) {
	p.Data[p.Index] = d
	p.Advance(1)
}

//ReadBytes is experimental and has not been tested
func (p *Packet)ReadBytes(amount int) []byte {
	return p.Data[p.Index:p.Index+p.Advance(amount)]
}
