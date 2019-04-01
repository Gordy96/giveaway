package account

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type AccountStatus string

const (
	New         AccountStatus = "new"
	Available   AccountStatus = "available"
	Busy        AccountStatus = "busy"
	Maintenance AccountStatus = "maintenance"
	CheckPoint  AccountStatus = "checkpoint"
	Error       AccountStatus = "error"
)

type WorksWithAccount interface {
	GetAccount() *Account
	SetAccount(*Account)
}

type AuthenticatesAccount interface {
	WorksWithAccount
	Login() (bool, error)
}

type Account struct {
	Id       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`

	DeviceId string `json:"device_id" bson:"device_id"`
	PhoneId  string `json:"phone_id" bson:"phone_id"`
	AdId     string `json:"adid" bson:"adid"`
	GUID     string `json:"guid" bson:"guid"`

	Cookies map[string][]*http.Cookie `json:"cookies" bson:"cookies"`
	Proxy   string                    `json:"proxy" bson:"proxy"`

	Status AccountStatus `json:"status" bson:"status"`

	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
	CreatedAt int64 `json:"created_at" bson:"created_at"`
}

func GenerateDeviceId() string {
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomizer.Seed(time.Now().UnixNano())
	buf := make([]byte, 8)
	randomizer.Read(buf)
	return fmt.Sprintf("android-%s", hex.EncodeToString(buf))
}

func NewAccount(username, password string) *Account {
	ac := &Account{}
	ac.Username = username
	ac.Password = password

	ac.AdId = uuid.New().String()
	ac.DeviceId = GenerateDeviceId()
	ac.PhoneId = uuid.New().String()
	ac.GUID = uuid.New().String()

	ac.Cookies = make(map[string][]*http.Cookie)

	ac.CreatedAt = time.Now().UnixNano()
	ac.UpdatedAt = time.Now().UnixNano()

	return ac
}

func (a *Account) MakeCookieJar() http.CookieJar {
	jar, _ := cookiejar.New(nil)
	for u, c := range a.Cookies {
		uri, _ := url.Parse(u)
		jar.SetCookies(uri, c)
	}

	return jar
}

func (a *Account) SetCookies(path *url.URL, cookies []*http.Cookie) {
	a.Cookies[path.String()] = cookies
}

func (a *Account) ResetCookies() {
	a.Cookies = make(map[string][]*http.Cookie)
}

func NewPredefinedAccount(username, password, id, adId, deviceId, phoneId, guid string) *Account {
	ac := &Account{}
	ac.Username = username
	ac.Password = password

	ac.Id = id
	ac.AdId = adId
	ac.DeviceId = deviceId
	ac.PhoneId = phoneId
	ac.GUID = guid

	ac.Cookies = make(map[string][]*http.Cookie)

	ac.CreatedAt = time.Now().Unix()
	ac.UpdatedAt = time.Now().Unix()

	return ac
}
