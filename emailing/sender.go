package emailing

import "sync/atomic"

type Sender interface {
	Send(email *Email) error
}

type wellKnownEmailSenders struct {
	LoungeUpContact   string
	LoungeUpNoReply   string
	LoungeUpOPS       string
	LoungeUpSales     string
	LoungeUpSupport   string
	LoungeUpTranslate string
}

var defaultWellKnownEmailSenders atomic.Pointer[wellKnownEmailSenders]

//nolint:gochecknoinits
func init() {
	defaultWellKnownEmailSenders.Store(&wellKnownEmailSenders{
		LoungeUpContact:   "contact@loungeup.com",
		LoungeUpNoReply:   "noreply@loungeup.com",
		LoungeUpOPS:       "ops@loungeup.com",
		LoungeUpSales:     "sales@loungeup.com",
		LoungeUpSupport:   "support@loungeup.com",
		LoungeUpTranslate: "translate@loungeup.com",
	})
}

func WellKnownEmailSenders() *wellKnownEmailSenders { return defaultWellKnownEmailSenders.Load() }
