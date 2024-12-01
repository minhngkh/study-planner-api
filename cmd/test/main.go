package main

import "fmt"

type testInterface interface {
	incrementID()
}

type testStruct struct {
	id int
}

func (t testStruct) incrementID() {
	t.id++
}

func testMethod(i testInterface) {
	i.incrementID()
}

func main() {
	var t = testStruct{id: 1}

	testMethod(t)

	fmt.Println(t.id)
}
