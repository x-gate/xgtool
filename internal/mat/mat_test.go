package mat

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewMatrix(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}

	mat, err := NewMatrix(data, 2, 3)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(mat, Matrix{2, 3, []int{1, 2, 3, 4, 5, 6}}); diff != "" {
		t.Errorf("NewMatrix() mismatch (-want +got):\n%s", diff)
	}
}

func TestMatrix_Rotate(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}

	mat, _ := NewMatrix(data, 2, 3)

	if diff := cmp.Diff(mat.Rotate(), Matrix{3, 2, []int{2, 4, 6, 1, 3, 5}}); diff != "" {
		t.Errorf("Matrix.Rotate() mismatch (-want +got):\n%s", diff)
	}
}
