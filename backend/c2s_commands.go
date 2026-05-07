package main

type C2SCommand string

const (
	C2SAddItem           C2SCommand = "addItem"
	C2SDeleteItem        C2SCommand = "deleteItem"
	C2SRestoreItem       C2SCommand = "restoreItem"
	C2SSetItemCompletion C2SCommand = "setItemCompletion"
)
