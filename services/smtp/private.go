package smtp

import (
	"bytes"
	"net/mail"
)

func encodeAddrs(addrs []mail.Address, bufTo *bytes.Buffer) {
	for i := range addrs {
		if i != 0 {
			bufTo.Write([]byte(", "))
		}
		bufTo.WriteString(addrs[i].String())
	}
}
