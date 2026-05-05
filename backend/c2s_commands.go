package main

type C2SCommand string

const (
	C2SAddItem           C2SCommand = "addItem"
	C2SDeleteItem        C2SCommand = "deleteItem"
	C2SSetItemCompletion C2SCommand = "setItemCompletion"
)
