package mat

import (
	"errors"
	"fmt"
)

// Matrix a 1D array that represents a 2D array.
type Matrix struct {
	W    int
	H    int
	Data []int
}

// ErrInvalidData is returned when the data length is not equal to (width * height)
var ErrInvalidData = errors.New("invalid data")

// NewMatrix creates a new Matrix from the given data, width and height.
func NewMatrix[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](data []T, w, h int) (mat Matrix, err error) {
	if len(data) != w*h {
		return mat, fmt.Errorf("%w: len(data)=%d, w=%d, h=%d", ErrInvalidData, len(data), w, h)
	}

	mat = Matrix{
		W:    w,
		H:    h,
		Data: make([]int, w*h),
	}

	for i := 0; i < w*h; i++ {
		mat.Data[i] = int(data[i])
	}

	return
}

// Rotate rotates the matrix -90 degrees clockwise.
func (m Matrix) Rotate() Matrix {
	dst := Matrix{
		W:    m.H,
		H:    m.W,
		Data: make([]int, m.W*m.H),
	}
	col, row := m.W, m.H

	for i := 0; i < m.H*m.W; i++ {
		// Because Matrix.Data is a 1D array, the dst[i] is calculated.
		//
		// There's a 3x4 matrix:		Rotate -90 degrees:
		// 	[							[
		// 		[1, 2, 3],					[3, 6, 9, 12],
		// 		[4, 5, 6],					[2, 5, 8, 11],
		// 		[7, 8, 9],					[1, 4, 7, 10],
		// 		[10, 11, 12],			]
		//	]
		//
		// col, row := 3, 4
		// for i := 0; i < col*row; i+=row {
		// 		dst[i+0] = src[col*1-1-i/row]
		// 		dst[i+1] = src[col*2-1-i/row]
		// 		dst[i+2] = src[col*3-1-i/row]
		// 		dst[i+3] = src[col*4-1-i/row]
		// }
		dst.Data[i] = m.Data[col*(i%row+1)-1-i/row]
	}
	return dst
}

// ToUint16Array converts the matrix to a uint16 array.
//
// Note: the matrix data is converted to uint16, so the data may be lost.
func (m Matrix) ToUint16Array() (arr []uint16) {
	arr = make([]uint16, m.W*m.H)

	for i := 0; i < m.W*m.H; i++ {
		arr[i] = uint16(m.Data[i])
	}

	return
}
