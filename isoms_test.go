package caplib

import "fmt"

func ExampleSpace_Translations() {
	for d := 1; d <= 2; d++ {
		s := NewSpace(d)
		fmt.Println(s.Translations().Perms)
	}

	// Output:
	// [[0 1 2] [1 2 0] [2 0 1]]
	// [[0 1 2 3 4 5 6 7 8] [1 2 0 4 5 3 7 8 6] [2 0 1 5 3 4 8 6 7] [3 4 5 6 7 8 0 1 2] [4 5 3 7 8 6 1 2 0] [5 3 4 8 6 7 2 0 1] [6 7 8 0 1 2 3 4 5] [7 8 6 1 2 0 4 5 3] [8 6 7 2 0 1 5 3 4]]
}

func ExampleSpace_CoordPerms() {
	for d := 1; d <= 3; d++ {
		s := NewSpace(d)
		fmt.Println(s.CoordPerms().Perms)
	}

	// Output:
	// [[0 1 2]]
	// [[0 1 2 3 4 5 6 7 8] [0 3 6 1 4 7 2 5 8]]
	// [[0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26] [0 3 6 1 4 7 2 5 8 9 12 15 10 13 16 11 14 17 18 21 24 19 22 25 20 23 26] [0 1 2 9 10 11 18 19 20 3 4 5 12 13 14 21 22 23 6 7 8 15 16 17 24 25 26] [0 3 6 9 12 15 18 21 24 1 4 7 10 13 16 19 22 25 2 5 8 11 14 17 20 23 26] [0 9 18 1 10 19 2 11 20 3 12 21 4 13 22 5 14 23 6 15 24 7 16 25 8 17 26] [0 9 18 3 12 21 6 15 24 1 10 19 4 13 22 7 16 25 2 11 20 5 14 23 8 17 26]]
}

func ExampleSpace_CoordReflections() {
	for d := 1; d <= 2; d++ {
		s := NewSpace(d)
		fmt.Println(s.CoordReflections().Perms)
	}

	// Output:
	// [[0 1 2] [0 2 1]]
	// [[0 1 2 3 4 5 6 7 8] [0 2 1 3 5 4 6 8 7] [0 1 2 6 7 8 3 4 5] [0 2 1 6 8 7 3 5 4]]
}

func ExampleSpace_LinearIsoms() {
	for d := 1; d <= 4; d++ {
		s := NewSpace(d)
		isoms := s.LinearIsoms()
		fmt.Println(isoms.Len())
		if d < 3 {
			fmt.Println(isoms.Perms1.Perms)
			fmt.Println(isoms.Perms2.Perms)
		} else {
			fmt.Println(isoms.Perms1.Len())
			fmt.Println(isoms.Perms2.Len())
		}

	}

	// Output:
	// 2
	// [[0 1 2] [0 2 1]]
	// [[0 1 2]]
	// 48
	// [[0 1 2 3 4 5 6 7 8] [0 2 1 3 5 4 6 8 7] [0 1 2 6 7 8 3 4 5] [0 2 1 6 8 7 3 5 4] [0 3 6 1 4 7 2 5 8] [0 3 6 2 5 8 1 4 7] [0 6 3 1 7 4 2 8 5] [0 6 3 2 8 5 1 7 4]]
	// [[0 1 2 3 4 5 6 7 8] [0 1 2 4 5 3 8 6 7] [0 1 2 5 3 4 7 8 6] [0 3 6 4 7 1 8 2 5] [0 3 6 5 8 2 7 1 4] [0 4 8 5 6 1 7 2 3]]
	// 11232
	// 48
	// 234
	// 24261120
	// 384
	// 63180
}

func ExampleSpace_LinearIsomsModCoords() {
	fmt.Println(NewSpace(1).LinearIsomsModCoords().Perms)
	fmt.Println(NewSpace(2).LinearIsomsModCoords().Perms)
	fmt.Println(len(NewSpace(3).LinearIsomsModCoords().Perms))

	// Output:
	// [[0 1 2]]
	// [[0 1 2 3 4 5 6 7 8] [0 1 2 4 5 3 8 6 7] [0 1 2 5 3 4 7 8 6] [0 3 6 4 7 1 8 2 5] [0 3 6 5 8 2 7 1 4] [0 4 8 5 6 1 7 2 3]]
	// 234
}

func ExampleSpace_LinearIsomsFixingCounts() {
	fmt.Println(NewSpace(1).LinearIsomsFixingCounts([]int{0, 1, 2}).Perms)
	fmt.Println(NewSpace(1).LinearIsomsFixingCounts([]int{0, 1, 1}).Perms)
	fmt.Println(NewSpace(1).LinearIsomsFixingCounts([]int{0, 0, 0}).Perms)

	fmt.Println(NewSpace(2).LinearIsomsFixingCounts([]int{0, 1, 1, 2, 0, 0, 2, 0, 0}).Perms)
	fmt.Println(NewSpace(2).LinearIsomsFixingCounts([]int{0, 1, 1, 1, 0, 0, 1, 0, 0}).Perms)
	fmt.Println(NewSpace(2).LinearIsomsFixingCounts([]int{0, 0, 0, 0, 0, 0, 0, 0, 0}).Perms)

	// Output:
	// [[0 1 2]]
	// [[0 1 2] [0 2 1]]
	// [[0 1 2] [0 2 1]]
	// [[0 1 2 3 4 5 6 7 8] [0 1 2 6 7 8 3 4 5] [0 2 1 3 5 4 6 8 7] [0 2 1 6 8 7 3 5 4]]
	// [[0 1 2 3 4 5 6 7 8] [0 1 2 6 7 8 3 4 5] [0 2 1 3 5 4 6 8 7] [0 2 1 6 8 7 3 5 4] [0 3 6 1 4 7 2 5 8] [0 3 6 2 5 8 1 4 7] [0 6 3 1 7 4 2 8 5] [0 6 3 2 8 5 1 7 4]]
	// [[0 1 2 3 4 5 6 7 8] [0 1 2 4 5 3 8 6 7] [0 1 2 5 3 4 7 8 6] [0 1 2 6 7 8 3 4 5] [0 1 2 7 8 6 5 3 4] [0 1 2 8 6 7 4 5 3] [0 2 1 3 5 4 6 8 7] [0 2 1 4 3 5 8 7 6] [0 2 1 5 4 3 7 6 8] [0 2 1 6 8 7 3 5 4] [0 2 1 7 6 8 5 4 3] [0 2 1 8 7 6 4 3 5] [0 3 6 1 4 7 2 5 8] [0 3 6 2 5 8 1 4 7] [0 3 6 4 7 1 8 2 5] [0 3 6 5 8 2 7 1 4] [0 3 6 7 1 4 5 8 2] [0 3 6 8 2 5 4 7 1] [0 4 8 1 5 6 2 3 7] [0 4 8 2 3 7 1 5 6] [0 4 8 3 7 2 6 1 5] [0 4 8 5 6 1 7 2 3] [0 4 8 6 1 5 3 7 2] [0 4 8 7 2 3 5 6 1] [0 5 7 1 3 8 2 4 6] [0 5 7 2 4 6 1 3 8] [0 5 7 3 8 1 6 2 4] [0 5 7 4 6 2 8 1 3] [0 5 7 6 2 4 3 8 1] [0 5 7 8 1 3 4 6 2] [0 6 3 1 7 4 2 8 5] [0 6 3 2 8 5 1 7 4] [0 6 3 4 1 7 8 5 2] [0 6 3 5 2 8 7 4 1] [0 6 3 7 4 1 5 2 8] [0 6 3 8 5 2 4 1 7] [0 7 5 1 8 3 2 6 4] [0 7 5 2 6 4 1 8 3] [0 7 5 3 1 8 6 4 2] [0 7 5 4 2 6 8 3 1] [0 7 5 6 4 2 3 1 8] [0 7 5 8 3 1 4 2 6] [0 8 4 1 6 5 2 7 3] [0 8 4 2 7 3 1 6 5] [0 8 4 3 2 7 6 5 1] [0 8 4 5 1 6 7 3 2] [0 8 4 6 5 1 3 2 7] [0 8 4 7 3 2 5 1 6]]
}
