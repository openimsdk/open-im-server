package utilsv2

import (
	"fmt"
	"testing"
)

func TestDistinct(t *testing.T) {
	arr := []int{1, 1, 1, 4, 4, 5, 2, 3, 3, 3, 6}
	fmt.Println(Distinct(arr))
}

func TestDeleteAt(t *testing.T) {
	arr := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(DeleteAt(arr, 0, 1, -1, -2))
	fmt.Println(DeleteAt(arr))
	fmt.Println(DeleteAt(arr, 1))
}

func TestSliceToMap(t *testing.T) {
	type Item struct {
		ID   string
		Name string
	}
	list := []Item{
		{ID: "111", Name: "111"},
		{ID: "222", Name: "222"},
		{ID: "333", Name: "333"},
	}

	m := SliceToMap(list, func(t Item) string {
		return t.ID
	})

	fmt.Printf("%+v\n", m)

}

func TestIndexOf(t *testing.T) {
	arr := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	fmt.Println(IndexOf(arr, 3))

}
