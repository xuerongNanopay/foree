package scotia

import (
	"testing"
)

// var testPrivateKey = `
// -----BEGIN RSA PRIVATE KEY-----
// MIIBPQIBAAJBALqbHeRLCyOdykC5SDLqI49ArYGYG1mqaH9/GnWjGavZM02fos4l
// c2w6tCchcUBNtJvGqKwhC5JEnx3RYoSX2ucCAwEAAQJBAKn6O+tFFDt4MtBsNcDz
// GDsYDjQbCubNW+yvKbn4PJ0UZoEebwmvH1ouKaUuacJcsiQkKzTHleu4krYGUGO1
// mEECIQD0dUhj71vb1rN1pmTOhQOGB9GN1mygcxaIFOWW8znLRwIhAMNqlfLijUs6
// rY+h1pJa/3Fh1HTSOCCCCWA0NRFnMANhAiEAwddKGqxPO6goz26s2rHQlHQYr47K
// vgPkZu2jDCo7trsCIQC/PSfRsnSkEqCX18GtKPCjfSH10WSsK5YRWAY3KcyLAQIh
// AL70wdUu5jMm2ex5cZGkZLRB50yE6rBiHCd5W1WdTFoe
// -----END RSA PRIVATE KEY-----
// `

// var testPublicKey = `
// ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAID6DY9+sMypB1GacOpkKY9Uk36hG1LfzTaKUF7A+oKR5 test@xrw.io

// `

func TestSignJwt(t *testing.T) {
	//TODO: This is only a example of how to sign jwt from private key string.
	t.Run("SignJwt should work", func(t *testing.T) {
		// encoded := `eyJhbGciOiJSUzI1NiIsImtpZCI6ImtrayIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ4dWUiLCJzdWIiOiJ4dWUiLCJhdWQiOlsiZm9vIl0sImV4cCI6MTcyMjQ4MTc3MX0.IfmawCozC7gLltFo9XfvHN4TNPs7OLW5tlNHY0B7of1-Hv2P6qVwXP8ELXkeAlENgGPZdPa1Z011k6We99oWaA    `
		// block, _ := pem.Decode([]byte(testPrivateKey))
		// key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
		// // key := parseResult.(*rsa.PrivateKey)

		// client := &ScotiaClientImpl{
		// 	Config: ScotiaConfig{
		// 		ClientId:    "xue",
		// 		JWTKid:      "kkk",
		// 		JWTAudience: "foo",
		// 		JWTExpiry:   5,
		// 		JWTSignKey:  key,
		// 	},
		// }
		// ss, err := client.signJWT()
		// if err != nil {
		// 	t.Errorf("signing fail with error %v", err.Error())
		//  }
		// fmt.Println(ss)
		// fmt.Println(encoded)
		// if ss != encoded {
		// 	t.Errorf("signing fail")
		// }

	})
}
