package main

import (
	"fmt"
	"github.com/hyperledger/xcoin/proxy"
)

func main() {
	run := proxy.NewAppRunner()
	args := []string{"adduser",
		"{\"pubkey\":\"MIIBCgKCAQEA1QVZzTjW/uq6tk/Oct6rO/HyPZB5+xSt2n7bJy7XPFRNbfeKV0gGYeHJy8ctbDCN8TxU2evxrPr5QaXcJsOHBoJdo9MLaxwirz5bT2Ctom7W2hIfUVcafzTvbRtpAZCkS+ZwGjn/u3/gsJqF0HUHZlmobGL9JxF0BF8vqD/x7VD0qaPlwPQNRk3cyuywIR/a1kAjxiXjWUrmnFvqpNZwO5P0+/KwQ3D8v/PS+s1ZaG1SaFqPSHm9CZp7wknin3/0prLPxnVmtUv3lsKP6lv4Vc6i3OSlBWBIlSIXG9GjDgCcAHvDZnixgoa6jv7Zif2qYfOTqwL/rTSgWPGRm3oWCwIDAQAB\",\"timestamp\":100}", "signature"}
	text := run.SendRequest(args)
	fmt.Println(text)
}
