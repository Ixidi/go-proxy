package conn

import (
	"bytes"
	"errors"
	"reflect"
)

type Packet interface {
	Id() VarInt
	IdInt() int
	Data() []byte
	UnpackFields(fields ...interface{ DataType }) error
	Unpack(v interface{}) error
}

type packet struct {
	id   VarInt
	data []byte
}

func (p *packet) Id() VarInt {
	return p.id
}

func (p *packet) IdInt() int {
	return int(p.id)
}

func (p *packet) Data() []byte {
	return p.data
}

func (p *packet) UnpackFields(fields ...interface{ DataType }) error {
	reader := bytes.NewReader(p.Data())
	for _, f := range fields {
		if err := f.Read(reader); err != nil {
			return err
		}
	}

	return nil
}

func (p *packet) Unpack(s interface{}) error {
	if reflect.TypeOf(s).Kind() != reflect.Ptr {
		return errors.New("pointer to struct expected")
	}

	reader := bytes.NewReader(p.Data())
	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		dataType := f.Interface()

		var data interface{}
		var err error

		//TODO better solution syf
		switch x := dataType.(type) {
		case VarInt:
			err = x.Read(reader)
			data = x
		case UShort:
			err = x.Read(reader)
			data = x
		case String:
			err = x.Read(reader)
			data = x
		case Long:
			err = x.Read(reader)
			data = x
		case UUID:
			err = x.Read(reader)
			data = x
		case ByteArray:
			err = x.Read(reader)
			data = x
		case Position:
			err = x.Read(reader)
			data = x
		case Float:
			err = x.Read(reader)
			data = x
		case Double:
			err = x.Read(reader)
			data = x
		case Byte:
			err = x.Read(reader)
			data = x
		case UByte:
			err = x.Read(reader)
			data = x
		case Bool:
			err = x.Read(reader)
			data = x
		case Int:
			err = x.Read(reader)
			data = x
		case StringArray:
			err = x.Read(reader)
			data = x
		case NBT:
			err = x.Read(reader)
			data = x
		}

		if err != nil {
			return err
		}

		f.Set(reflect.ValueOf(data))
	}

	return nil
}

func Pack(id VarInt, s interface{}) (Packet, error) {
	if reflect.TypeOf(s).Kind() != reflect.Ptr {
		return nil, errors.New("pointer to struct expected")
	}

	var buff bytes.Buffer
	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		dataType := f.Interface()

		var err error

		//TODO better solution
		switch x := dataType.(type) {
		case VarInt:
			err = x.Write(&buff)
		case UShort:
			err = x.Write(&buff)
		case String:
			err = x.Write(&buff)
		case Long:
			err = x.Write(&buff)
		case UUID:
			err = x.Write(&buff)
		case ByteArray:
			err = x.Write(&buff)
		case Position:
			err = x.Write(&buff)
		case Float:
			err = x.Write(&buff)
		case Double:
			err = x.Write(&buff)
		case Byte:
			err = x.Write(&buff)
		case UByte:
			err = x.Write(&buff)
		case Bool:
			err = x.Write(&buff)
		case Int:
			err = x.Write(&buff)
		case StringArray:
			err = x.Write(&buff)
		case NBT:
			err = x.Write(&buff)
		}

		if err != nil {
			return nil, err
		}
	}

	return NewPacket(id, buff.Bytes()), nil
}

func NewPacket(id VarInt, data []byte) Packet {
	return &packet{
		id:   id,
		data: data,
	}
}

func CreatePacket(id VarInt, fields ...interface{ DataType }) (Packet, error) {
	var buff bytes.Buffer
	for _, f := range fields {
		if err := f.Read(&buff); err != nil {
			return nil, err
		}
	}

	return NewPacket(id, buff.Bytes()), nil
}
