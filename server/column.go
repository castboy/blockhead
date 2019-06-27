package server

import (
	"trading-system/common/decimal/errors"
)

type Col struct {
	ID       int
	position int
	length   int
}

type Cols struct {
	allCols []*Col
}

const SeparatorLength int = 1

func (c *Cols) Get(colID int) (*Col, error) {
	if err := c.ColIDCheck(colID); err != nil {
		return nil, err
	}

	if colID == -1 { // last col.
		return c.allCols[len(c.allCols)-1], nil
	}

	return c.allCols[colID], nil
}

func (c *Cols) ColIDCheck(colID int) error {
	if colID < -1 || colID > len(c.allCols)-1 {
		return errors.New("invalid col id")
	}

	return nil
}

func (c *Cols) Rm(colID int) error {
	if err := c.ColIDCheck(colID); err != nil {
		return err
	}

	if colID == -1 { // last col.
		c.allCols = c.allCols[:len(c.allCols)-1]
		return nil
	}

	for i := range c.allCols {
		if i > colID {
			c.allCols[i].ID--
			c.allCols[i].position -= c.allCols[colID].length + SeparatorLength
		}
	}
	c.allCols = append(c.allCols[:colID], c.allCols[colID+1:]...)

	return nil
}

func (c *Cols) InsertBefore(colID int, col *Col) error {
	if err := c.ColIDCheck(colID); err != nil {
		return err
	}

	if colID == -1 { // insert after last col.
		c.allCols = append(c.allCols, col)
		return nil
	}

	tail := make([]*Col, 0)

	for i := range c.allCols {
		if i >= colID {
			c.allCols[i].ID++
			c.allCols[i].position += c.allCols[colID].length + SeparatorLength

			tail = append(tail, c.allCols[i])
		}
	}

	c.allCols = append(c.allCols[:colID], col)
	c.allCols = append(c.allCols, tail...)

	return nil
}

func (c *Cols) Exchange(col1ID, col2ID int) error {
	if col1ID == col2ID {
		return errors.New("can not be same id")
	}

	if err := c.ColIDCheck(col1ID); err != nil {
		return err
	}
	if err := c.ColIDCheck(col2ID); err != nil {
		return err
	}

	if col1ID == - 1 {
		col1ID = len(c.allCols) - 1
	}

	if col2ID == - 1 {
		col2ID = len(c.allCols) - 1
	}

	col1, col2 := c.allCols[col1ID], c.allCols[col2ID]
	col1.ID, col2.ID = col2ID, col1ID

	c.Rm(col1ID)
	c.Rm(col2ID)

	c.InsertBefore(col1ID, col2)
	c.InsertBefore(col2ID, col1)

	return nil
}

func (c *Cols) Combine(col1ID, col2ID int) error {
	if err := c.ColIDCheck(col1ID); err != nil {
		return err
	}
	if err := c.ColIDCheck(col2ID); err != nil {
		return err
	}

	if col1ID == - 1 {
		col1ID = len(c.allCols) - 1
	}

	if col2ID == - 1 {
		col2ID = len(c.allCols) - 1
	}

	if col1ID == col2ID {
		return errors.New("can not be same id")
	}

	smallerID, biggerID := col1ID, col2ID

	if smallerID > biggerID {
		smallerID, biggerID = col2ID, col1ID
	}

	if biggerID - smallerID > 1 {
		return errors.New("not adjacent two cols.")
	}

	for i := range c.allCols {
		if i == smallerID {
			c.allCols[smallerID].length += c.allCols[biggerID].length
		}
		if i > biggerID {
			c.allCols[i].ID--
			c.allCols[i].position -= SeparatorLength //how length of separator. TODO
		}
	}

	if biggerID == len(c.allCols) {
		c.allCols = c.allCols[:biggerID]
		return nil
	}

	c.allCols = append(c.allCols[:biggerID], c.allCols[biggerID+1:]...)
	return nil
}

func (c *Cols) Split(colID int, position uint) error { // position is the second col start position.
	if err := c.ColIDCheck(colID); err != nil {
		return err
	}

	if colID == -1 {
		colID = len(c.allCols)-1
	}

	newCol := &Col{}
	for i := range c.allCols {
		if i == colID {
			oldLength := c.allCols[colID].length
			c.allCols[colID].length = int(position)
			newCol.ID = colID + 1
			newCol.position = c.allCols[colID].position+c.allCols[colID].length+SeparatorLength
			newCol.length = oldLength - c.allCols[colID].length
		}

		if i > colID {
			c.allCols[i].ID++
			c.allCols[i].position += SeparatorLength
		}
	}

	if colID == len(c.allCols)-1 {
		c.allCols = append(c.allCols, newCol)
		return nil
	}

	head, tail := c.allCols[:colID+1], c.allCols[colID+1:]

	c.allCols = append(head, newCol)
	c.allCols = append(c.allCols, tail...)

	return nil
}