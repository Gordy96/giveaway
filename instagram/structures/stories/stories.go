package stories

import (
	"giveaway/data/owner"
	"strconv"
)

type Candidate struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	URL    string  `json:"url"`
}
type ImageVersions2 struct {
	Candidates []Candidate `json:"candidates"`
}
type VideoVersion struct {
	Type   int     `json:"type"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	URL    string  `json:"url"`
	ID     string  `json:"id"`
}
type ReelMention struct { //TODO: Mention of account
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        float64 `json:"z"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Rotation float64 `json:"rotation"`
	IsPinned int     `json:"is_pinned"`
	IsHidden int     `json:"is_hidden"`
	User     User    `json:"user"`
}
type StoryFeedMedia struct { //TODO: Re-share of post
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Z           float64 `json:"z"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	Rotation    float64 `json:"rotation"`
	IsPinned    int     `json:"is_pinned"`
	IsHidden    int     `json:"is_hidden"`
	MediaID     int64   `json:"media_id"`
	ProductType string  `json:"product_type"`
}
type CreativeConfig struct {
	FaceEffectID                int64    `json:"face_effect_id"`
	PersistedEffectMetadataJSON string   `json:"persisted_effect_metadata_json,omitempty"`
	CameraFacing                string   `json:"camera_facing"`
	ShouldRenderTryItOn         bool     `json:"should_render_try_it_on"`
	EffectID                    int64    `json:"effect_id"`
	Name                        string   `json:"name"`
	EffectName                  string   `json:"effect_name"`
	EffectActions               []string `json:"effect_actions"`
	AttributionUsername         string   `json:"attribution_username"`
	AttributionID               int      `json:"attribution_id"`
	CaptureType                 string   `json:"capture_type,omitempty"`
}
type Link struct {
	LinkType                                int         `json:"linkType"`
	WebURI                                  string      `json:"webUri"`
	AndroidClass                            string      `json:"androidClass"`
	Package                                 string      `json:"package"`
	DeeplinkURI                             string      `json:"deeplinkUri"`
	CallToActionTitle                       string      `json:"callToActionTitle"`
	RedirectURI                             interface{} `json:"redirectUri"`
	LeadGenFormID                           string      `json:"leadGenFormId"`
	IgUserID                                string      `json:"igUserId"`
	AppInstallObjectiveInvalidationBehavior interface{} `json:"appInstallObjectiveInvalidationBehavior"`
}
type StoryCta struct { //TODO: external links
	Links []Link `json:"links"`
}
type Hashtag struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}
type StoryHashtag struct { //TODO: hashtags
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        float64 `json:"z"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Rotation float64 `json:"rotation"`
	IsPinned int     `json:"is_pinned"`
	IsHidden int     `json:"is_hidden"`
	Hashtag  Hashtag `json:"hashtag"`
}
type Tallies struct {
	Text     string  `json:"text"`
	FontSize float64 `json:"font_size"`
	Count    int     `json:"count"`
}
type PollSticker struct {
	ID               string      `json:"id"`
	PollID           int64       `json:"poll_id"`
	Question         string      `json:"question"`
	Tallies          []Tallies   `json:"tallies"`
	PromotionTallies interface{} `json:"promotion_tallies"`
	ViewerCanVote    bool        `json:"viewer_can_vote"`
	IsSharedResult   bool        `json:"is_shared_result"`
	Finished         bool        `json:"finished"`
}
type StoryPoll struct { //TODO: polls
	X           float64     `json:"x"`
	Y           float64     `json:"y"`
	Z           float64     `json:"z"`
	Width       float64     `json:"width"`
	Height      float64     `json:"height"`
	Rotation    float64     `json:"rotation"`
	IsPinned    int         `json:"is_pinned"`
	IsHidden    int         `json:"is_hidden"`
	PollSticker PollSticker `json:"poll_sticker"`
}
type Location struct {
	Pk               int     `json:"pk"`
	Name             string  `json:"name"`
	Address          string  `json:"address"`
	City             string  `json:"city"`
	ShortName        string  `json:"short_name"`
	Lng              float64 `json:"lng"`
	Lat              float64 `json:"lat"`
	ExternalSource   string  `json:"external_source"`
	FacebookPlacesID int64   `json:"facebook_places_id"`
}
type StoryLocation struct { //TODO: locations
	X        float64  `json:"x"`
	Y        float64  `json:"y"`
	Z        float64  `json:"z"`
	Width    float64  `json:"width"`
	Height   float64  `json:"height"`
	Rotation float64  `json:"rotation"`
	IsPinned int      `json:"is_pinned"`
	IsHidden int      `json:"is_hidden"`
	Location Location `json:"location"`
}
type StoryItem struct {
	TakenAt                  int64            `json:"taken_at"`
	Pk                       int64            `json:"pk"`
	ID                       string           `json:"id"`
	DeviceTimestamp          int64            `json:"device_timestamp"`
	MediaType                int              `json:"media_type"`
	Code                     string           `json:"code"`
	ClientCacheKey           string           `json:"client_cache_key"`
	FilterType               int              `json:"filter_type"`
	ImageVersions2           ImageVersions2   `json:"image_versions2"`
	OriginalWidth            float64          `json:"original_width"`
	OriginalHeight           float64          `json:"original_height"`
	CaptionPosition          float64          `json:"caption_position"`
	IsReelMedia              bool             `json:"is_reel_media"`
	IsDashEligible           int              `json:"is_dash_eligible,omitempty"`
	VideoDashManifest        string           `json:"video_dash_manifest,omitempty"`
	VideoCodec               string           `json:"video_codec,omitempty"`
	NumberOfQualities        int              `json:"number_of_qualities,omitempty"`
	HasAudio                 bool             `json:"has_audio,omitempty"`
	VideoDuration            float64          `json:"video_duration,omitempty"`
	User                     User             `json:"user"`
	Caption                  interface{}      `json:"caption"`
	CaptionIsEdited          bool             `json:"caption_is_edited"`
	PhotoOfYou               bool             `json:"photo_of_you"`
	CanViewerSave            bool             `json:"can_viewer_save"`
	OrganicTrackingToken     string           `json:"organic_tracking_token"`
	ExpiringAt               int64            `json:"expiring_at"`
	ImportedTakenAt          int64            `json:"imported_taken_at,omitempty"`
	CanReshare               bool             `json:"can_reshare"`
	CanReply                 bool             `json:"can_reply"`
	SupportsReelReactions    bool             `json:"supports_reel_reactions"`
	ShowOneTapFbShareTooltip bool             `json:"show_one_tap_fb_share_tooltip"`
	HasSharedToFb            int              `json:"has_shared_to_fb"`
	CreativeConfig           *CreativeConfig  `json:"creative_config,omitempty"`
	AdAction                 string           `json:"ad_action,omitempty"`
	LinkText                 string           `json:"link_text,omitempty"`
	ReelMentions             []ReelMention    `json:"reel_mentions,omitempty"`
	StoryFeedMedia           []StoryFeedMedia `json:"story_feed_media,omitempty"`
	VideoVersions            []VideoVersion   `json:"video_versions,omitempty"`
	StoryCta                 []StoryCta       `json:"story_cta,omitempty"`
	StoryHashtags            []StoryHashtag   `json:"story_hashtags,omitempty"`
	StoryPolls               []StoryPoll      `json:"story_polls,omitempty"`
	StoryLocations           []StoryLocation  `json:"story_locations,omitempty"`
}

func (s *StoryItem) GetKey() interface{} {
	return strconv.FormatInt(s.User.Pk, 10)
}

func (s *StoryItem) GetCreationDate() int64 {
	return s.TakenAt
}

func (s *StoryItem) GetOwner() *owner.Owner {
	return &owner.Owner{
		Id:       strconv.FormatInt(s.User.Pk, 10),
		Username: s.User.Username,
	}
}

type User struct {
	Pk                         int64  `json:"pk"`
	Username                   string `json:"username"`
	FullName                   string `json:"full_name"`
	IsPrivate                  bool   `json:"is_private"`
	ProfilePicURL              string `json:"profile_pic_url"`
	ProfilePicID               string `json:"profile_pic_id,omitempty"`
	IsVerified                 bool   `json:"is_verified"`
	HasAnonymousProfilePicture bool   `json:"has_anonymous_profile_picture,omitempty"`
	IsUnpublished              bool   `json:"is_unpublished,omitempty"`
	IsFavorite                 bool   `json:"is_favorite,omitempty"`
}
