package gopdf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unicode/utf16"
)

type outlinesObj struct {
	buffer     bytes.Buffer
	getRoot    func() *GoPdf
	FirstId    int
	LastId     int
	innerIndex int
}

func (ol *outlinesObj) init(funcGetRoot func() *GoPdf) {
	ol.getRoot = funcGetRoot
}

func (ol *outlinesObj) build(objID int) error {
	ol.buffer.WriteString("<< /Type /" + ol.getType())
	ol.buffer.WriteString(fmt.Sprintf(" /First %d 0 R /Last %d 0 R ", ol.FirstId, ol.LastId))
	ol.buffer.WriteString(">>\n")
	return nil
}

func (ol *outlinesObj) getType() string {
	return "Outlines"
}

func (ol *outlinesObj) getObjBuff() *bytes.Buffer {
	return &(ol.buffer)
}

type OutlineObj struct {
	ParentId int
	PrevId   int
	NextId   int
	FirstId  int
	LastId   int
	PageId   int
	ActionId int
	Title    string
	buffer   bytes.Buffer
}

func convertUTF16ToBigEndianBytes(arr []uint16) []byte {
	buf := make([]byte, 2*len(arr))
	for idx, val := range arr {
		binary.BigEndian.PutUint16(buf[idx*2:], val)
	}
	return buf
}

func (ol *OutlineObj) init(funcGetRoot func() *GoPdf) {}
func (ol *OutlineObj) build(objID int) error {
	ol.buffer.WriteString("<< /Title (")
	u16Arr := utf16.Encode([]rune(ol.Title))
	ol.buffer.Write([]byte{0xfe, 0xff})
	ol.buffer.Write(convertUTF16ToBigEndianBytes(u16Arr))
	ol.buffer.WriteString(") ")
	if ol.ParentId > 0 {
		ol.buffer.WriteString(fmt.Sprintf("/Parent %d 0 R ", ol.ParentId))
	}
	if ol.PrevId > 0 {
		ol.buffer.WriteString(fmt.Sprintf("/Prev %d 0 R ", ol.PrevId))
	}
	if ol.NextId > 0 {
		ol.buffer.WriteString(fmt.Sprintf("/Next %d 0 R ", ol.NextId))
	}
	if ol.FirstId > 0 {
		ol.buffer.WriteString(fmt.Sprintf("/First %d 0 R ", ol.FirstId))
	}
	if ol.LastId > 0 {
		ol.buffer.WriteString(fmt.Sprintf("/Last %d 0 R ", ol.LastId))
	}
	if ol.ActionId > 0 {
		ol.buffer.WriteString(fmt.Sprintf("/A %d 0 R ", ol.ActionId))
	} else {
		ol.buffer.WriteString(fmt.Sprintf("/Dest [ %d 0 R ] ", ol.PageId))
	}
	ol.buffer.WriteString("/F 0 >>\n")
	return nil
}
func (ol *OutlineObj) getType() string {
	return "OutlineItem"
}
func (ol *OutlineObj) getObjBuff() *bytes.Buffer {
	return &(ol.buffer)
}

type actionObj struct {
	buffer  bytes.Buffer
	PageId  int
	YOffset float64
}

func (a *actionObj) init(funcGetRoot func() *GoPdf) {
}

func (a *actionObj) build(objID int) error {
	a.buffer.WriteString("<< /S /GoTo /Type /" + a.getType())
	a.buffer.WriteString(fmt.Sprintf(" /D [%d 0 R /XYZ null %f 0.0] >>\n", a.PageId, a.YOffset))
	return nil
}

func (a *actionObj) getType() string {
	return "Action"
}

func (a *actionObj) getObjBuff() *bytes.Buffer {
	return &(a.buffer)
}
