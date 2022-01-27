package converter

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"

	"dav_converter/pkg/dav/exception"
	"dav_converter/pkg/dav/usecase"
)

// GetHeaderDav получаем заголовок DAV
func getHeaderDav(f io.Reader) (usecase.HeaderDav, error) {
	header := usecase.HeaderDav{}
	size := 32
	b := make([]byte, size)

	countReader, err := f.Read(b)

	if err != nil {
		return header, err
	}

	if countReader != size {
		return header, exception.FewBytes{}
	}

	initStruct(&header, b)

	return header, nil
}

func getHeaderFrame(f io.Reader) (usecase.FrameHeader, error) {
	header := usecase.FrameHeader{}
	size := 32
	b := make([]byte, size)

	countReader, err := f.Read(b)

	if err != nil {
		return header, err
	}

	if countReader != size {
		return header, exception.FewBytes{}
	}

	initStruct(&header, b)

	return header, nil
}

func getTagIdx(f io.ReadSeeker) (usecase.TagFrameIDX, error) {
	size := 20
	tag := usecase.TagFrameIDX{}
	b := make([]byte, size)

	countReader, err := f.Read(b)

	if err != nil {
		return tag, err
	}

	if countReader != size {
		return tag, exception.FewBytes{}
	}

	initStruct(&tag, b)

	return tag, nil
}

func initStruct(items interface{}, data []byte) {
	r := bytes.NewReader(data)

	val := reflect.ValueOf(items).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		switch typeField.Type.Kind() {
		case reflect.Int64:
			var value int64
			binary.Read(r, binary.LittleEndian, &value)
			valueField.Set(reflect.ValueOf(value))

		case reflect.Int32:
			var value int32
			binary.Read(r, binary.LittleEndian, &value)
			valueField.Set(reflect.ValueOf(value))

		case reflect.Int16:
			var value int16
			binary.Read(r, binary.LittleEndian, &value)
			valueField.Set(reflect.ValueOf(value))

		case reflect.Int8:
			var value int8
			binary.Read(r, binary.LittleEndian, &value)
			valueField.Set(reflect.ValueOf(value))

		case reflect.Uint16:
			var value uint16
			binary.Read(r, binary.LittleEndian, &value)
			valueField.Set(reflect.ValueOf(value))

		case reflect.Uint32:
			var value uint32
			binary.Read(r, binary.LittleEndian, &value)
			valueField.Set(reflect.ValueOf(value))
		}
	}
}
