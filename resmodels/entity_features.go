package resmodels

type EntityFeatures struct {
	Newsletters               EntityFeature
	InstantMessages           EntityFeature
	Events                    EntityFeature
	Emailing                  EntityFeature
	GuestProfile              EntityFeature
	ArrivalsManagement        EntityFeature
	Application               EntityFeature
	NewslettersBySendinblue   EntityFeature
	GroupEmailingCampaigns    EntityFeature
	GroupApplicationCampaigns EntityFeature
}

type EntityFeature struct {
	Data      map[string]any `json:"data"`
	Activated bool           `json:"activated"`
}

type RawEntityFeature struct {
	FeatureName string         `json:"name"`
	Activated   bool           `json:"activated"`
	JsonData    map[string]any `json:"data"`
}

type FeatureName string

const (
	NewsLettersName               FeatureName = "newsletters"
	InstantMessagesName           FeatureName = "instantMessages"
	EventsName                    FeatureName = "events"
	EmailingName                  FeatureName = "emailing"
	GuestProfileName              FeatureName = "guestProfile"
	ArrivalsManagementName        FeatureName = "arrivalsManagement"
	ApplicationName               FeatureName = "application"
	NewslettersBySendinblueName   FeatureName = "newslettersBySendinblue"
	GroupEmailingCampaignsName    FeatureName = "groupEmailingCampaigns"
	GroupApplicationCampaignsName FeatureName = "groupApplicationCampaigns"
)

func (f FeatureName) String() string { return string(f) }

func MapRawEntityFeaturesToEntityFeatures(rawEntityFeatures []*RawEntityFeature) *EntityFeatures {
	result := &EntityFeatures{}

	for _, rawEntityFeature := range rawEntityFeatures {
		feature := EntityFeature{
			Data:      rawEntityFeature.JsonData,
			Activated: rawEntityFeature.Activated,
		}

		switch rawEntityFeature.FeatureName {
		case NewsLettersName.String():
			result.Newsletters = feature
		case InstantMessagesName.String():
			result.InstantMessages = feature
		case EventsName.String():
			result.Events = feature
		case EmailingName.String():
			result.Emailing = feature
		case GuestProfileName.String():
			result.GuestProfile = feature
		case ArrivalsManagementName.String():
			result.ArrivalsManagement = feature
		case ApplicationName.String():
			result.Application = feature
		case NewslettersBySendinblueName.String():
			result.NewslettersBySendinblue = feature
		case GroupEmailingCampaignsName.String():
			result.GroupEmailingCampaigns = feature
		case GroupApplicationCampaignsName.String():
			result.GroupApplicationCampaigns = feature
		}
	}

	return result
}
