package naming

import (
	"crypto/rsa"
	"fmt"

	c "../common"
)

type NamingImpl struct {
	lookupTable    map[string]*c.AOR
	keylookupTable map[string]*rsa.PublicKey
}

func NewNamingImpl() *NamingImpl {
	return &NamingImpl{lookupTable: map[string]*c.AOR{}}
}

func (n *NamingImpl) Lookup(topic string) *c.AOR {
	fmt.Println("looking up")
	aor := n.lookupTable[topic]
	fmt.Println(aor)
	return aor
}

func (n *NamingImpl) Register(topic string, aor *c.AOR) {
	fmt.Println("Registering: " + aor.ToString())
	n.lookupTable[topic] = aor
}

func (n *NamingImpl) RegisterKey(topic string, key *rsa.PublicKey) {
	n.keylookupTable[topic] = key
}

func (n *NamingImpl) LookupKey(topic string) *rsa.PublicKey {
	key := n.keylookupTable[topic]
	return key
}
