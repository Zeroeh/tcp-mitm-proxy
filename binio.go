package main

import (
	"encoding/binary"
	"fmt"
)

const (
	packetSize = 20000
)

func (p *Packet) Advance(amount int) int {
	p.Index += amount
	return amount
}

type Packet struct {
	Index  int
	Length int
	ID     int
	Data   []byte //holds ALL of our data
}

func NewPacket() Packet {
	p := Packet{}
	p.Init()
	
	return p
}
func (p *Packet) Init() {
	p.Length = 0
	p.ID = 0
	p.Index = 0 //start at the first byte of data
	p.Data = make([]byte, packetSize)
}

func (p *Packet)ReadString() string {
	length := int(p.ReadUInt16()) //absolute
	if length == 0 {
		return ""
	}
	var str []byte
	str = p.Data[p.Index:p.Index+length]
	fmt.Println("Str data is", str)
	p.Advance(length)
	return string(str)
}

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
	length := int(p.ReadUInt32())
	if length == 0 {
		return ""
	}
	var str []byte
	str = p.Data[p.Index:p.Index+length]
	p.Advance(length)
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

//the float functions need to be optimized...
func (p *Packet)ReadFloat() float32 {
	floatBuf := p.Data[p.Index:p.Index+p.Advance(4)]
	newData := make([]byte, 4)
	newData[0] = floatBuf[3]
	newData[1] = floatBuf[2]
	newData[2] = floatBuf[1]
	newData[3] = floatBuf[0]
	return float32(binary.BigEndian.Uint32(newData))
}

func (p *Packet)WriteFloat(f float32) {
	p.WriteUInt32(uint32(f))
}

func (p *Packet)ReadInt16() int16 {
	return int16(binary.BigEndian.Uint16(p.Data[p.Index:p.Index+2]))
}

func (p *Packet)WriteInt16(i int16) {
	binary.BigEndian.PutUint16(p.Data[p.Index:p.Index+p.Advance(2)], uint16(i))
}

func (p *Packet)ReadUInt16() uint16 {
	return binary.BigEndian.Uint16(p.Data[p.Index:p.Index+2])
}

func (p *Packet)WriteUInt16(i uint16) {
	binary.BigEndian.PutUint16(p.Data[p.Index:p.Index+p.Advance(2)], i)
}

func (p *Packet)ReadInt32() int32 {
	return int32(binary.BigEndian.Uint32(p.Data[p.Index:p.Index+4]))
}

func (p *Packet)WriteInt32(i int32) {
	binary.BigEndian.PutUint32(p.Data[p.Index:p.Index+p.Advance(4)], uint32(i))
}

func (p *Packet)ReadUInt32() uint32 {
	return binary.BigEndian.Uint32(p.Data[p.Index:p.Index+4])
}

func (p *Packet)WriteUInt32(i uint32) {
	binary.BigEndian.PutUint32(p.Data[p.Index:p.Index+p.Advance(4)], i)
}

func (p *Packet)ReadByte() byte {
	return p.Data[p.Index:p.Index+p.Advance(1)][0]
}

func (p *Packet)WriteByte(d byte) {
	p.Data[p.Index] = d
	p.Advance(1)
}

func (p *Packet)ReadBytes(amount int) []byte {
	return p.Data[p.Index:p.Index+p.Advance(amount)]
}

func AppendData(slice []byte, data ...byte) []byte {
    m := len(slice)
    n := m + len(data)
    if n > cap(slice) { // if necessary, reallocate
        //allocate double what's needed, for future growth.
        newSlice := make([]byte, (n+1)*2)
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[0:n]
    copy(slice[m:n], data)
    return slice
}
