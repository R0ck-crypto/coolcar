package token

import (
	"github.com/golang-jwt/jwt"
	"testing"
	"time"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzau8ZWtrCV4pit2F5zCL
iyQKHut7nx9ilgCKQ1PtHwvHuup8wSBvXar/XQHLRuAzAM0sJmsEKCfHh1uij4ez
QwLuXkvPJlmFA4tEYD7xWX+XW/A/++vWVrhvKiLzzJp7oB9Eizxe56mz0GzLiAnI
IoD95GkEcpb4DoN7vi+N9lhGKU+iS5w1BUy5ENJgLmbHDQ6i6VkCli8DQYi/YZy6
eBuVahmuZjwLvZcx9IMfnbqHzwlywPioKgODsJmck4KP89pCuPGzGj7W8Gc88IZy
U2HvwkHtXshIpaIqOs6mjCy77XxlDosACFHVywxpOPBIs6i0aBa11uLs6z064E3+
+QIDAQAB
-----END PUBLIC KEY-----`

func TestJWTVerifier_Verify(t *testing.T) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		t.Fatalf("cannot parse public key: %v", err)
	}

	j := &JWTVerifier{
		PublicKey: pubKey,
	}

	cases := []struct {
		name    string
		tkn     string
		now     time.Time
		want    string
		wantErr bool
	}{
		{
			name:    "valid_token",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDQyMzcwODAsImlhdCI6MTY0NDIyOTg4MCwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjFmZTUyN2E4YjQ4MmI1NDBjN2JjM2Y0In0.ejBNhK55sBXPPkPMBeRZEDNjfCIvEMpPVc3XUTOZEjC_KQhCnFEL3NKAWJEek10E83S6-LdpB8BHsRIE5IJNTYwfYYGpAPLyk3WkV-NeP1AWctmK763hL3dlwXTjtuZbW00z7XFy4kELz3pPov8DRahjYsxviA-gq2fciZWC2MYatl0-as5Seq3SNtCQV7PF_OOC1M8iixoX5uohtksGpdtztLl2rFumwCHWk0lj-vJFQ2Am3aBjQ9lMLZbnfWpykDtOw4bzRVvqGMnYBVFTdoKgkH_zW-Eyg4gm8ofvVowvvYZ_F5VsN9odN-hpv0jTzr88pu_819xSvuyJs59xqA",
			now:     time.Unix(1644235080, 0),
			want:    "61fe527a8b482b540c7bc3f4",
			wantErr: false,
		},
		{
			name:    "token_expired",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDQyMzcwODAsImlhdCI6MTY0NDIyOTg4MCwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjFmZTUyN2E4YjQ4MmI1NDBjN2JjM2Y0In0.ejBNhK55sBXPPkPMBeRZEDNjfCIvEMpPVc3XUTOZEjC_KQhCnFEL3NKAWJEek10E83S6-LdpB8BHsRIE5IJNTYwfYYGpAPLyk3WkV-NeP1AWctmK763hL3dlwXTjtuZbW00z7XFy4kELz3pPov8DRahjYsxviA-gq2fciZWC2MYatl0-as5Seq3SNtCQV7PF_OOC1M8iixoX5uohtksGpdtztLl2rFumwCHWk0lj-vJFQ2Am3aBjQ9lMLZbnfWpykDtOw4bzRVvqGMnYBVFTdoKgkH_zW-Eyg4gm8ofvVowvvYZ_F5VsN9odN-hpv0jTzr88pu_819xSvuyJs59xqA",
			now:     time.Unix(1644239080, 0),
			want:    "61fe527a8b482b540c7bc3f4",
			wantErr: true,
		},
		{
			name:    "bad_token",
			tkn:     "JhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDQyMzcwODAsImlhdCI6MTY0NDIyOTg4MCwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjFmZTUyN2E4YjQ4MmI1NDBjN2JjM2Y0In0.ejBNhK55sBXPPkPMBeRZEDNjfCIvEMpPVc3XUTOZEjC_KQhCnFEL3NKAWJEek10E83S6-LdpB8BHsRIE5IJNTYwfYYGpAPLyk3WkV-NeP1AWctmK763hL3dlwXTjtuZbW00z7XFy4kELz3pPov8DRahjYsxviA-gq2fciZWC2MYatl0-as5Seq3SNtCQV7PF_OOC1M8iixoX5uohtksGpdtztLl2rFumwCHWk0lj-vJFQ2Am3aBjQ9lMLZbnfWpykDtOw4bzRVvqGMnYBVFTdoKgkH_zW-Eyg4gm8ofvVowvvYZ_F5VsN9odN-hpv0jTzr88pu_819xSvuyJs59xqA",
			now:     time.Unix(1644235080, 0),
			want:    "61fe527a8b482b540c7bc3f4",
			wantErr: true,
		},
		{
			name:    "wrong_signature",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ8.eyJleHAiOjE2NDQyMzcwODAsImlhdCI6MTY0NDIyOTg4MCwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjFmZTUyN2E4YjQ4MmI1NDBjN2JjM2Y0In0.ejBNhK55sBXPPkPMBeRZEDNjfCIvEMpPVc3XUTOZEjC_KQhCnFEL3NKAWJEek10E83S6-LdpB8BHsRIE5IJNTYwfYYGpAPLyk3WkV-NeP1AWctmK763hL3dlwXTjtuZbW00z7XFy4kELz3pPov8DRahjYsxviA-gq2fciZWC2MYatl0-as5Seq3SNtCQV7PF_OOC1M8iixoX5uohtksGpdtztLl2rFumwCHWk0lj-vJFQ2Am3aBjQ9lMLZbnfWpykDtOw4bzRVvqGMnYBVFTdoKgkH_zW-Eyg4gm8ofvVowvvYZ_F5VsN9odN-hpv0jTzr88pu_819xSvuyJs59xqA",
			now:     time.Unix(1644235080, 0),
			want:    "61fe527a8b482b540c7bc3f4",
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			jwt.TimeFunc = func() time.Time {
				return c.now
			}
			accountID, err := j.Verify(c.tkn)

			if !c.wantErr && err != nil {
				t.Errorf("verification failed:%v", err)
			}
			if c.wantErr && err == nil {
				t.Errorf("want error;got no error")
			}
			if accountID != c.want {
				t.Errorf("wrong account id. want:%q, got:%q", c.want, accountID)
			}
		})
	}
}
