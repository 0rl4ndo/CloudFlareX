package main

import (
	"CloudFlareX/internal/CF"
	"fmt"
)

func madin() {
	S := cloudflare.GetSession("")
	Params, Js, Err := S.GetCfParams()

	if Err != nil {
		return
	}

	Cook := S.GetCookie(Params, Js)
	
	IsFlagged := true
	if len(Cook) == 189 {
		IsFlagged = false
	}

	fmt.Printf(`
 InvisibleJs:
    - EncryptionKey: %s
    - Nonce: %s
    - S: %s
	
 CfParams:
	- R: %s
	- M: %s
	- U: %s
	- S: %v
 
  Output:
    - Flagged: %t
    - Cookie: %s
`, Js.Password, Js.Nonce, Js.S, Params.R, Params.M, Params.U, Params.S, IsFlagged, Cook)
}