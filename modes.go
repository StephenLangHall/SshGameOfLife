package main

import (
	"os"
	"log"
	"fmt"
	"strconv"
)

func SaveBoard(b [h][w]bool) {
	d1 := []byte{}
	fmt.Println("Hight -> "+strconv.Itoa(h)+"\nWidth -> "+strconv.Itoa(w)+"\n")
	d1 = append(d1, byte(h), byte(w))
	for y, row := range b {
		for x, cell := range row {
			if cell {
				fmt.Println(strconv.Itoa(x) + " x " + strconv.Itoa(y))
				d1 = append(d1, byte(y))
				d1 = append(d1, byte(x))
			}
		}
	}
	err := os.WriteFile("./statesave", d1, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("===================")
}

func OpenBoard() [h][w]bool {
	nb := [h][w]bool{}
	nb = [len(nb)][len(nb[0])]bool{}

	dat, err := os.ReadFile("./statesave")
	if err != nil {
		log.Fatal(err)
	}
	if len(dat) < 3 {
		return nb
	}
	nh := int(dat[1])
	nw := int(dat[0])
	if h != nh || w != nw {
		log.Fatal("Hight "+strconv.Itoa(h)+" -> "+strconv.Itoa(nh)+"\nWidth "+strconv.Itoa(w)+" -> "+strconv.Itoa(nw)+"\n")
	}
	data := dat[2:]
	for i := range len(data)/2 {
		fmt.Println(i)
		ty := int(data[i*2])
		tx := int(data[i*2+1])
		fmt.Println(strconv.Itoa(tx) + " x " + strconv.Itoa(ty))
		if ty < nh && tx < nw {
			nb[ty][tx] = true
		}
	}
	fmt.Println("Hight "+strconv.Itoa(h)+" -> "+strconv.Itoa(nh)+"\nWidth "+strconv.Itoa(w)+" -> "+strconv.Itoa(nw))
	fmt.Println("===================")
	return nb
}

