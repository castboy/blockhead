package main

import (
	"blockhead/server"
	"fmt"
)

type Blockheader struct {
	FullTextOperater
	ColumnOperater
}

// full text.
type FullTextOperater interface {
	//Analyzer
}

// column.
type ColumnOperater interface {
	Delete(col int)
	MoveAhead(col int)
	MoveAfter(col int)
}


// word.
type WordOperater interface {
	WordOperate(atom *server.WordOpAtom)
}

func main() {
	w := server.Word{
		Ori: []rune("zz"),
	}

	strs := []string{"f", "af", "f", "fe", "fde", "fe"}
	L := len(strs)
	for i := range strs {
		if i + 1 < L {
			w.Analize([]rune(strs[i]), []rune(strs[i+1]))
		}
	}

	for i := range w.Ops {
		w.WordOperate(w.Ops[i])
		fmt.Println(string(w.Ori))
	}
}