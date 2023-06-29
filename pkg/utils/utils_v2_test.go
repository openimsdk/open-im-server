package utils

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
	fmt.Println(Delete(arr, 0, 1, -1, -2))
	fmt.Println(Delete(arr))
	fmt.Println(Delete(arr, 1))
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

	fmt.Println(IndexOf(3, arr...))

}

func TestSort(t *testing.T) {
	arr := []int{1, 1, 1, 4, 4, 5, 2, 3, 3, 3, 6}
	fmt.Println(Sort(arr, false))
}

func TestBothExist(t *testing.T) {
	arr1 := []int{1, 1, 1, 4, 4, 5, 2, 3, 3, 3, 6}
	arr2 := []int{6, 1, 3}
	arr3 := []int{5, 1, 3, 6}
	fmt.Println(BothExist(arr1, arr2, arr3))
}

func TestCompleteAny(t *testing.T) {
	type Item struct {
		ID    int
		Value string
	}

	ids := []int{1, 2, 3, 4, 5, 6, 7, 8}

	var list []Item

	for _, id := range ids {
		list = append(list, Item{
			ID:    id,
			Value: fmt.Sprintf("%d", id*1000),
		})
	}

	DeleteAt(&list, -1)
	DeleteAt(&ids, -1)

	ok := Complete(ids, Slice(list, func(t Item) int {
		return t.ID
	}))

	fmt.Printf("%+v\n", ok)

}
