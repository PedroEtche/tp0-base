package common

import (
	"strconv"
)

type Bet struct {
	ID        uint8
	Name      string
	LastName  string
	Document  uint32
	Birthdate string //"YYYY-MM-DD"
	Number    uint32
}

func CreatBetFromCSVLine(agency uint8, line []string) (*Bet, error) {
	document, err := strconv.ParseUint(line[2], 10, 32)
	if err != nil {
		return nil, err
	}
	number, err := strconv.ParseUint(line[4], 10, 32)
	if err != nil {
		return nil, err
	}
	return &Bet{
		ID:        agency,
		Name:      line[0],
		LastName:  line[1],
		Document:  uint32(document),
		Birthdate: line[3],
		Number:    uint32(number),
	}, nil
}

func (b *Bet) GetYear() uint16 {
	year, _ := strconv.ParseUint(b.Birthdate[0:4], 10, 16)
	return uint16(year)
}

func (b *Bet) GetMonth() uint8 {
	month, _ := strconv.ParseUint(b.Birthdate[5:7], 10, 8)
	return uint8(month)
}

func (b *Bet) GetDay() uint8 {
	day, _ := strconv.ParseUint(b.Birthdate[8:10], 10, 8)
	return uint8(day)
}
