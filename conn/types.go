package conn

import (
	"encoding/binary"
	"errors"
	"github.com/Tnze/go-mc/nbt"
	"github.com/satori/go.uuid"
	"io"
)

type DataType interface {
	Read(reader io.Reader) error
	Write(writer io.Writer) error
}

type VarInt int32
type UShort uint16
type String string
type Float float32
type Double float64
type Byte int8
type UByte uint8
type Bool bool
type Int int32
type Long int64
type UUID uuid.UUID
type ByteArray []byte
type Position struct {
	X int
	Y int
	Z int
}
type StringArray []String
type NBT struct {
	V interface{}
}

func (v *VarInt) Read(reader io.Reader) error {
	numRead := 0
	result := 0
	var read byte
	var err error
	for {
		bytes := make([]byte, 1)
		if _, err = reader.Read(bytes); err != nil {
			return err
		}
		read = bytes[0]
		value := read & 0b01111111
		result |= int(value) << (7 * numRead)

		numRead++
		if numRead > 5 {
			return errors.New("VarInt is too big")
		}

		if read&0b10000000 == 0 {
			break
		}
	}

	*v = VarInt(result)
	return nil
}

func (v VarInt) Write(writer io.Writer) error {
	value := v
	for {
		temp := byte(value & 0b01111111)
		value >>= 7
		if value != 0 {
			temp |= 0b10000000
		}
		if _, err := writer.Write([]byte{temp}); err != nil {
			return err
		}
		if value == 0 {
			break
		}
	}
	return nil
}

func (v *UShort) Read(reader io.Reader) error {
	return binary.Read(reader, binary.BigEndian, v)
}

func (v UShort) Write(writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, v)
}

func (v *String) Read(reader io.Reader) error {
	var l VarInt
	var err error

	if err = l.Read(reader); err != nil {
		return err
	}

	buff := make([]byte, l)
	if _, err = reader.Read(buff); err != nil {
		return err
	}

	*v = String(buff)
	return nil
}

func (v String) Write(writer io.Writer) error {
	var err error

	buff := []byte(v)
	if err = VarInt(len(buff)).Write(writer); err != nil {
		return err
	}

	if _, err = writer.Write(buff); err != nil {
		return err
	}

	return nil
}

func (v *Long) Read(reader io.Reader) error {
	return binary.Read(reader, binary.BigEndian, v)
}

func (v Long) Write(writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, v)
}

func (u *UUID) Read(reader io.Reader) error {
	buff := make([]byte, 16)
	if _, err := reader.Read(buff); err != nil {
		return err
	}
	v, err := uuid.FromBytes(buff)
	if err != nil {
		return err
	}
	*u = UUID(v)
	return nil
}

func (u UUID) Write(writer io.Writer) error {
	v := uuid.UUID(u)
	if _, err := writer.Write(v.Bytes()); err != nil {
		return err
	}
	return nil
}

func (b *ByteArray) Read(reader io.Reader) error {
	var length VarInt
	if err := length.Read(reader); err != nil {
		return err
	}
	buff := make([]byte, length)
	if _, err := reader.Read(buff); err != nil {
		return err
	}

	*b = buff
	return nil
}

func (b ByteArray) Write(writer io.Writer) error {
	if err := VarInt(len(b)).Write(writer); err != nil {
		return err
	}

	if _, err := writer.Write(b); err != nil {
		return err
	}

	return nil
}

func (p *Position) Read(reader io.Reader) error {
	var long Long
	err := long.Read(reader)
	if err != nil {
		return err
	}

	x := long >> 38
	y := long & 0xfff
	z := long << 26 >> 38
	if x >= 2^25 {
		x -= 2 ^ 26
	}
	if y >= 2^11 {
		y -= 2 ^ 12
	}
	if z >= 2^25 {
		z -= 2 ^ 26
	}

	*p = Position{
		X: int(x),
		Y: int(y),
		Z: int(z),
	}

	return nil
}

func (p Position) Write(writer io.Writer) error {
	return Long(((p.X & 0x3FFFFFF) << 38) | ((p.Z & 0x3FFFFFF) << 12) | (p.Y & 0xFFF)).Write(writer)
}

func (f *Float) Read(reader io.Reader) error {
	return binary.Read(reader, binary.BigEndian, f)
}

func (f Float) Write(writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, f)
}

func (d *Double) Read(reader io.Reader) error {
	return binary.Read(reader, binary.BigEndian, d)
}

func (d Double) Write(writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, d)
}

func (b *Byte) Read(reader io.Reader) error {
	return binary.Read(reader, binary.BigEndian, b)
}

func (b Byte) Write(writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, b)
}

func (b *UByte) Read(reader io.Reader) error {
	return binary.Read(reader, binary.BigEndian, b)
}

func (b UByte) Write(writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, b)
}

func (b *Int) Read(reader io.Reader) error {
	return binary.Read(reader, binary.BigEndian, b)
}

func (b Int) Write(writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, b)
}

func (b *Bool) Read(reader io.Reader) error {
	buff := make([]byte, 1)
	if _, err := reader.Read(buff); err != nil {
		return err
	}

	*b = buff[0] == 1
	return nil
}

func (b Bool) Write(writer io.Writer) error {
	var by byte
	if b {
		by = 1
	} else {
		by = 0
	}
	return binary.Write(writer, binary.BigEndian, by)
}

func (s *StringArray) Read(reader io.Reader) error {
	var count VarInt
	err := count.Read(reader)
	if err != nil {
		return err
	}

	var array []String
	for i := 0; i < int(count); i++ {
		var str String
		err = str.Read(reader)
		if err != nil {
			return err
		}

		array = append(array, str)
	}

	*s = array
	return nil
}

func (s StringArray) Write(writer io.Writer) error {
	err := VarInt(len(s)).Write(writer)
	if err != nil {
		return err
	}

	for _, str := range s {
		err := str.Write(writer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (N *NBT) Read(reader io.Reader) error {
	return nbt.NewDecoder(reader).Decode(&N.V)
}

func (N NBT) Write(writer io.Writer) error {
	return nbt.NewEncoder(writer).Encode(N.V)
}
