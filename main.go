package main

import (
	cloudflare "CloudFlareX/internal/CloudFlare"
	"fmt"
)

func main() {
	for {
		cf := cloudflare.GetCfSession()
		invisible := cf.GetInvisible()
		Params := cf.CFCVParams()
		Invisible := cf.InvisibleJs(invisible)
		cf_bm := cf.GetResult(Invisible, Params)
		IsFlagged := true
		if len(cf_bm) == 125 {
			IsFlagged = false
		}
		println(fmt.Sprintf("🏁 %t |  %t", len(cf_bm), IsFlagged))
	}

}
