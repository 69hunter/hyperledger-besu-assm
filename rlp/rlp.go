package rlp

import "fmt"

const extraDataPrefix = "0xf87ea00000000000000000000000000000000000000000000000000000000000000000f854"
const extraDataPostfix = "808400000000c0"

type Adapter struct {
}

func (a *Adapter) EncodeFromAddresses(addresses []string) (string, error) {
	extraData := extraDataPrefix
	for _, addr := range addresses {
		extraData += fmt.Sprintf("94%s", addr)
	}
	extraData += extraDataPostfix
	return extraData, nil
}

func NewAdapter() *Adapter {
	return &Adapter{}
}
