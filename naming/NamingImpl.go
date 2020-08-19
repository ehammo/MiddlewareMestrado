package naming

import (
	c "../common"
)

type NamingImpl struct {
	lookupTable map[string]*c.AOR
}
func NewNamingImpl() *NamingImpl {
	return &NamingImpl{lookupTable: map[string]*c.AOR{}}
}

func (n *NamingImpl) lookup(topic string) *c.AOR {
	aor := n.lookupTable[topic]
	return aor
}

func (n *NamingImpl) register(topic string, aor *c.AOR) {
	n.lookupTable[topic] = aor
}



