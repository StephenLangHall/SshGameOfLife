package main

import (
	"os"
	"log"
	"fmt"
	"strconv"
)

func SaveBoard(b [h][w]bool) {
	d1 := []byte{}
	d1 = append(d1, byte(h), byte(w))
	for y, row := range b {
		for x, cell := range row {
			if cell {
				d1 = append(d1, byte(y), byte(x))
			}
		}
	}
	err := os.WriteFile("./statesave", d1, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func OpenBoard() [h][w]bool {
	nb := [h][w]bool{}
	nb = [len(nb)][len(nb[0])]bool{}

	dat, err := os.ReadFile("./statesave")
	if err != nil {
		log.Fatal(err)
	}
	nh := int(dat[1])
	nw := int(dat[0])
	if h != nh || w != nw {
		log.Fatal("Hight "+strconv.Itoa(h)+" -> "+strconv.Itoa(nh)+"\nWidth "+strconv.Itoa(w)+" -> "+strconv.Itoa(nw)+"\n")
	}
	data := dat[3:]
	for i := range len(data)/2 {
		ty := int(data[i*2])
		tx := int(data[i*2+1])
		fmt.Println(strconv.Itoa(tx) + " x " + strconv.Itoa(ty))
		if ty < nh && tx < nw {
			nb[ty][tx] = true
		}
	}
	fmt.Println("Hight "+strconv.Itoa(h)+" -> "+strconv.Itoa(nh)+"\nWidth "+strconv.Itoa(w)+" -> "+strconv.Itoa(nw))
	return nb
}

