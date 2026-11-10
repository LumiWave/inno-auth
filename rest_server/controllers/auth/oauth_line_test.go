package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// LINE ID 토큰 검증 API 흉내: 토큰마다 다른 사용자(sub), 다른 채널 ID·모르는 토큰은 400.
// client_id 를 안 보내면 통과시키는 최악의 경우를 가정 → 채널 ID 미설정 가드가 막아야 함.
func TestOauthLineVerifySocialKey(t *testing.T) {
	subs := map[string]string{"good": "U0123456789abcdef", "good2": "Ufedcba9876543210"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cid, sub := r.FormValue("client_id"), subs[r.FormValue("id_token")]
		if r.Method != http.MethodPost || (cid != "" && cid != "1234567890") || sub == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid_request","error_description":"Invalid IdToken."}`))
			return
		}
		w.Write([]byte(`{"iss":"https://access.line.me","sub":"` + sub + `","aud":"1234567890","email":"a@b.c"}`))
	}))
	defer srv.Close()
	old := LineVerifyIDTokenEndpoint
	LineVerifyIDTokenEndpoint = srv.URL
	defer func() { LineVerifyIDTokenEndpoint = old }()

	o := &OauthLine{SocialType: SocialType_LINE, ChannelID: "1234567890"}
	if id, ea, err := o.VerifySocialKey("good"); err != nil || id != "U0123456789abcdef" || ea != "a@b.c" {
		t.Fatalf("정상 토큰: id=%q ea=%q err=%v (sub 를 사용자 ID 로 써야 함)", id, ea, err)
	}
	if id, _, err := o.VerifySocialKey("good2"); err != nil || id != "Ufedcba9876543210" {
		t.Fatalf("다른 사용자가 다른 ID 를 받아야 함: id=%q err=%v", id, err)
	}
	if id, _, err := o.VerifySocialKey("bad"); err == nil || id != "" {
		t.Fatalf("LINE 이 거부한 토큰을 통과시킴: id=%q err=%v", id, err)
	}
	other := &OauthLine{SocialType: SocialType_LINE, ChannelID: "9999999999"}
	if id, _, err := other.VerifySocialKey("good"); err == nil || id != "" {
		t.Fatalf("다른 채널 ID 로 검증 통과: id=%q err=%v", id, err)
	}
	if _, _, err := (&OauthLine{SocialType: SocialType_LINE}).VerifySocialKey("good"); err == nil {
		t.Fatal("채널 ID 미설정인데 통과")
	}
}
