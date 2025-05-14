package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/LumiWave/baseutil/log"
)

type TelegramUser struct {
	UserID   int64  `json:"id"`
	UserName string `json:"username"`
}

type OauthTelegram struct {
	SocialType int64
}

func NewOauthTelegram() *OauthTelegram {
	return new(OauthTelegram)
}

func (o *OauthTelegram) GetSocialType() int64 {
	return o.SocialType
}

func (o *OauthTelegram) VerifySocialKey(initData string) (string, string, error) {
	// url decording
	decodedQuery, err := url.QueryUnescape(initData)
	if err != nil {
		log.Errorf("Error decoding URL: %v", err)
		return "", "", err
	}

	// 쿼리 파싱
	queries, err := url.ParseQuery(decodedQuery)
	if err != nil {
		log.Errorf("Error parsing query:%v", err)
		return "", "", err
	}
	// hash 정보만 따로 추출해서 verify에 사용
	hash := queries.Get("hash")
	queries.Del("hash")

	// 쿼리 파라미터를 '{key}={value}' 형식으로 변환하고 '\n'으로 연결
	var queryStrings []string
	for key, values := range queries {
		for _, value := range values {
			queryStrings = append(queryStrings, fmt.Sprintf("%s=%s", key, value))
		}
	}
	sort.Strings(queryStrings) // 정렬을 꼭 해야함
	formattedQueries := strings.Join(queryStrings, "\n")
	log.Debugf("Formatted Queries:%v", formattedQueries)

	// 시크릿 키 생성
	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(""))
	secretKey := h.Sum(nil)

	h = hmac.New(sha256.New, secretKey)
	h.Write([]byte(formattedQueries))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(calculatedHash), []byte(hash)) {
		log.Errorf("telegram verify fail cal:%v, origin:%v", calculatedHash, hash)
		return "", "", err
	}

	// 검증 성공 시 데이터 리턴
	user := new(TelegramUser)
	for k, v := range queries {
		if strings.EqualFold(k, "user") {
			json.Unmarshal([]byte(v[0]), &user)
			break
		}
	}

	return strconv.FormatInt(user.UserID, 10), user.UserName, nil
}
