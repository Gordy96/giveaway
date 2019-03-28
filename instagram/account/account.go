package account

import (
	"github.com/google/uuid"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Account struct {
	Id 			string							`json:"id" bson:"_id"`
	Username 	string							`json:"username" bson:"username"`
	Password 	string							`json:"password" bson:"password"`

	DeviceId	string							`json:"device_id" bson:"device_id"`
	PhoneId		string							`json:"phone_id" bson:"phone_id"`
	AdId		string							`json:"adid" bson:"adid"`
	GUID		string							`json:"guid" bson:"guid"`

	Cookies		map[string][]*http.Cookie		`json:"cookies" bson:"cookies"`
	Proxy		string
}

func NewAccount(username, password string) *Account {
	ac := &Account{}
	ac.Username = username
	ac.Password = password

	ac.AdId = uuid.New().String()
	ac.DeviceId = uuid.New().String()
	ac.PhoneId = uuid.New().String()
	ac.GUID = uuid.New().String()

	ac.Cookies = make(map[string][]*http.Cookie)

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

	return ac
}