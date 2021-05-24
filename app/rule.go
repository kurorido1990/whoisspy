package app

var spyNum = map[int]int{
	6:  1,
	9:  2,
	16: 3,
	20: 4,
}

func getSpyNum(num int) int {
	currentSpyNum := 1

	for playerNum, spyLimit := range spyNum {
		if playerNum <= num {
			currentSpyNum = spyLimit
			continue
		} else {
			return currentSpyNum
		}
	}

	return 5
}

var winNum = 3
