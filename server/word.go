package server

import "fmt"

type Word struct {
	Ori []rune
	Ops []*WordOpAtom
}

type WordOp int

const (
	PrefixRs WordOp = iota
	SuffixRs
	HeadTrimRs
	TailTrimRs
	MidRmRs
	MidInsertRs
)

type WordOpAtom struct {
	op       WordOp
	position int
	opNum 	 int
	runes    []rune
}

// analyzation

func (word *Word) Analize(ori, dst []rune) {
	if atom, ok := PrefixMode(ori, dst); ok {
		atom.op = PrefixRs
		word.Ops = append(word.Ops, atom)
	}

	if atom, ok := SuffixMode(ori, dst); ok {
		atom.op = SuffixRs
		word.Ops = append(word.Ops, atom)
	}

	if atom, ok := HeadTrimMode(ori, dst); ok {
		atom.op = HeadTrimRs
		word.Ops = append(word.Ops, atom)
	}

	if atom, ok := TailTrimMode(ori, dst); ok {
		atom.op = TailTrimRs
		word.Ops = append(word.Ops, atom)
	}

	if atom, ok := MidInsertMode(ori, dst); ok {
		atom.op = MidInsertRs
		word.Ops = append(word.Ops, atom)
	}

	if atom, ok := MidRmMode(ori, dst); ok {
		atom.op = MidRmRs
		word.Ops = append(word.Ops, atom)
	}
}

func PrefixMode(ori, dst []rune) (*WordOpAtom, bool) {
	oriL, dstL := len(ori), len(dst)

	if oriL >= dstL {
		return nil, false
	}

	spread := dstL - oriL

	i := dstL - 1
	for ; ; i-- {
		if i-spread < 0 {
			break
		}
		if dst[i] != ori[i-spread] {
			return nil, false
		}
	}

	return &WordOpAtom{
		runes: dst[:i+1],
	}, true
}

func SuffixMode(ori, dst []rune) (*WordOpAtom, bool) {
	oriL, dstL := len(ori), len(dst)

	if oriL >= dstL {
		return nil, false
	}

	i := 0
	for ; ; i++ {
		if i > oriL-1 {
			break
		}
		if dst[i] != ori[i] {
			return nil, false
		}
	}

	return &WordOpAtom{
		runes: dst[i:],
	}, true
}

func MidRmMode(ori, dst []rune) (*WordOpAtom, bool) {
	oriL, dstL := len(ori), len(dst)
	if oriL <= dstL {
		return nil, false
	}

	if ori[0] != dst[0] || ori[oriL-1] != dst[dstL-1] {
		return nil, false
	}

	i, spread := 0, oriL - dstL
	for ; ; i++ {
		if ori[i] == dst[i] && ori[oriL-1-i] == dst[dstL-1-i] {
			continue
		}
		if i > dstL-1-i {
			return &WordOpAtom{
				position: i,
				opNum: spread,
				runes: ori[i:i+spread],
			}, true
		}
	}

	return nil, false
}

func HeadTrimMode(ori, dst []rune) (*WordOpAtom, bool) {
	return PrefixMode(dst, ori)
}

func TailTrimMode(ori, dst []rune) (*WordOpAtom, bool) {
	return SuffixMode(dst, ori)
}

func MidInsertMode(ori, dst []rune) (*WordOpAtom, bool) {
	return MidRmMode(dst, ori)
}

// operation

func (word *Word) WordOperate(atom *WordOpAtom) {
	switch atom.op {
	case PrefixRs:
		word.Ori = append(atom.runes, word.Ori...)

	case SuffixRs:
		word.Ori = append(word.Ori, atom.runes...)

	case HeadTrimRs:
		word.Ori = word.Ori[len(atom.runes):]

	case TailTrimRs:
		word.Ori = word.Ori[:len(word.Ori)-len(atom.runes)]

	case MidRmRs:
		word.Ori = append(word.Ori[:atom.position], word.Ori[atom.position+atom.opNum:]...)
	case MidInsertRs:
		fmt.Println(word.Ori, "/////", atom)
		tail := word.Ori[atom.position:]
		word.Ori = append(word.Ori[:atom.position], atom.runes...)
		word.Ori = append(word.Ori, tail...)

	default:
		panic("invalid operation")
	}
}