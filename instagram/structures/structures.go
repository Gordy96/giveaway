package structures

import (
	"fmt"
	"giveaway/data"
	"strconv"
)

type PageInfo struct {
	HasNextPage bool   `json:"has_next_page"`
	EndCursor   string `json:"end_cursor"`
}

type ShortCodeMediaResponse struct {
	Data struct {
		ShortCodeMedia *struct {
			Typename   string `json:"__typename"`
			ID         string `json:"id"`
			ShortCode  string `json:"shortcode"`
			Dimensions struct {
				Height int `json:"height"`
				Width  int `json:"width"`
			} `json:"dimensions"`
			GatingInfo       interface{} `json:"gating_info"`
			MediaPreview     string      `json:"media_preview"`
			DisplayURL       string      `json:"display_url"`
			DisplayResources []struct {
				Src          string `json:"src"`
				ConfigWidth  int    `json:"config_width"`
				ConfigHeight int    `json:"config_height"`
			} `json:"display_resources"`
			AccessibilityCaption  string `json:"accessibility_caption"`
			IsVideo               bool   `json:"is_video"`
			ShouldLogClientEvent  bool   `json:"should_log_client_event"`
			TrackingToken         string `json:"tracking_token"`
			EdgeMediaToTaggedUser struct {
				Edges []interface{} `json:"edges"`
			} `json:"edge_media_to_tagged_user"`
			EdgeMediaToCaption struct {
				Edges []struct {
					Node struct {
						Text string `json:"text"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_media_to_caption"`
			CaptionIsEdited    bool `json:"caption_is_edited"`
			HasRankedComments  bool `json:"has_ranked_comments"`
			EdgeMediaToComment struct {
				Count    int64    `json:"count"`
				PageInfo PageInfo `json:"page_info"`
				Edges    []struct {
					Node struct {
						ID              string `json:"id"`
						Text            string `json:"text"`
						CreatedAt       int64  `json:"created_at"`
						DidReportAsSpam bool   `json:"did_report_as_spam"`
						Owner           struct {
							ID            string `json:"id"`
							IsVerified    bool   `json:"is_verified"`
							ProfilePicURL string `json:"profile_pic_url"`
							Username      string `json:"username"`
						} `json:"owner"`
						ViewerHasLiked bool `json:"viewer_has_liked"`
						EdgeLikedBy    struct {
							Count int64 `json:"count"`
						} `json:"edge_liked_by"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_media_to_comment"`
			CommentsDisabled     bool  `json:"comments_disabled"`
			TakenAtTimestamp     int64 `json:"taken_at_timestamp"`
			EdgeMediaPreviewLike struct {
				Count int64         `json:"count"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_media_preview_like"`
			EdgeMediaToSponsorUser struct {
				Edges []interface{} `json:"edges"`
			} `json:"edge_media_to_sponsor_user"`
			Location struct {
				ID            string `json:"id"`
				HasPublicPage bool   `json:"has_public_page"`
				Name          string `json:"name"`
				Slug          string `json:"slug"`
				AddressJSON   string `json:"address_json"`
			} `json:"location"`
			ViewerHasLiked             bool `json:"viewer_has_liked"`
			ViewerHasSaved             bool `json:"viewer_has_saved"`
			ViewerHasSavedToCollection bool `json:"viewer_has_saved_to_collection"`
			ViewerInPhotoOfYou         bool `json:"viewer_in_photo_of_you"`
			ViewerCanReShare           bool `json:"viewer_can_reshare"`
			Owner                      struct {
				ID                string `json:"id"`
				IsVerified        bool   `json:"is_verified"`
				ProfilePicURL     string `json:"profile_pic_url"`
				Username          string `json:"username"`
				BlockedByViewer   bool   `json:"blocked_by_viewer"`
				FollowedByViewer  bool   `json:"followed_by_viewer"`
				FullName          string `json:"full_name"`
				HasBlockedViewer  bool   `json:"has_blocked_viewer"`
				IsPrivate         bool   `json:"is_private"`
				IsUnpublished     bool   `json:"is_unpublished"`
				RequestedByViewer bool   `json:"requested_by_viewer"`
			} `json:"owner"`
			IsAd                       bool `json:"is_ad"`
			EdgeWebMediaToRelatedMedia struct {
				Edges []interface{} `json:"edges"`
			} `json:"edge_web_media_to_related_media"`
		} `json:"shortcode_media"`
	} `json:"data"`
	Status string `json:"status"`
}

type HashTagResponse struct {
	Data struct {
		HashTag *struct {
			ID                 string `json:"id"`
			Name               string `json:"name"`
			AllowFollowing     bool   `json:"allow_following"`
			IsFollowing        bool   `json:"is_following"`
			IsTopMediaOnly     bool   `json:"is_top_media_only"`
			ProfilePicURL      string `json:"profile_pic_url"`
			EdgeHashTagToMedia struct {
				Count    int64    `json:"count"`
				PageInfo PageInfo `json:"page_info"`
				Edges    []struct {
					Node struct {
						CommentsDisabled   bool   `json:"comments_disabled"`
						Typename           string `json:"__typename"`
						ID                 string `json:"id"`
						EdgeMediaToCaption struct {
							Edges []struct {
								Node struct {
									Text string `json:"text"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"edge_media_to_caption"`
						ShortCode          string `json:"shortcode"`
						EdgeMediaToComment struct {
							Count int64 `json:"count"`
						} `json:"edge_media_to_comment"`
						TakenAtTimestamp int64 `json:"taken_at_timestamp"`
						Dimensions       struct {
							Height int `json:"height"`
							Width  int `json:"width"`
						} `json:"dimensions"`
						DisplayURL  string `json:"display_url"`
						EdgeLikedBy struct {
							Count int64 `json:"count"`
						} `json:"edge_liked_by"`
						EdgeMediaPreviewLike struct {
							Count int64 `json:"count"`
						} `json:"edge_media_preview_like"`
						Owner struct {
							ID string `json:"id"`
						} `json:"owner"`
						ThumbnailSrc       string `json:"thumbnail_src"`
						ThumbnailResources []struct {
							Src          string `json:"src"`
							ConfigWidth  int    `json:"config_width"`
							ConfigHeight int    `json:"config_height"`
						} `json:"thumbnail_resources"`
						IsVideo              bool   `json:"is_video"`
						AccessibilityCaption string `json:"accessibility_caption"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_hashtag_to_media"`
			EdgeHashTagToTopPosts struct {
				Edges []struct {
					Node struct {
						Typename           string `json:"__typename"`
						ID                 string `json:"id"`
						EdgeMediaToCaption struct {
							Edges []struct {
								Node struct {
									Text string `json:"text"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"edge_media_to_caption"`
						ShortCode          string `json:"shortcode"`
						EdgeMediaToComment struct {
							Count int64 `json:"count"`
						} `json:"edge_media_to_comment"`
						TakenAtTimestamp int64 `json:"taken_at_timestamp"`
						Dimensions       struct {
							Height int `json:"height"`
							Width  int `json:"width"`
						} `json:"dimensions"`
						DisplayURL  string `json:"display_url"`
						EdgeLikedBy struct {
							Count int64 `json:"count"`
						} `json:"edge_liked_by"`
						EdgeMediaPreviewLike struct {
							Count int64 `json:"count"`
						} `json:"edge_media_preview_like"`
						Owner struct {
							ID string `json:"id"`
						} `json:"owner"`
						ThumbnailSrc       string `json:"thumbnail_src"`
						ThumbnailResources []struct {
							Src          string `json:"src"`
							ConfigWidth  int    `json:"config_width"`
							ConfigHeight int    `json:"config_height"`
						} `json:"thumbnail_resources"`
						IsVideo              bool   `json:"is_video"`
						AccessibilityCaption string `json:"accessibility_caption"`
					} `json:"node,omitempty"`
				} `json:"edges"`
			} `json:"edge_hashtag_to_top_posts"`
			EdgeHashTagToContentAdvisory struct {
				Count int64         `json:"count"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_hashtag_to_content_advisory"`
		} `json:"hashtag"`
	} `json:"data"`
	Status string `json:"status"`
}

type UserInfoResponse struct {
	LoggingPageID         string `json:"logging_page_id"`
	ShowSuggestedProfiles bool   `json:"show_suggested_profiles"`
	GraphQL               struct {
		User *struct {
			Biography              string `json:"biography"`
			BlockedByViewer        bool   `json:"blocked_by_viewer"`
			CountryBlock           bool   `json:"country_block"`
			ExternalURL            string `json:"external_url"`
			ExternalURLLinkShimmed string `json:"external_url_linkshimmed"`
			EdgeFollowedBy         struct {
				Count int64 `json:"count"`
			} `json:"edge_followed_by"`
			FollowedByViewer bool `json:"followed_by_viewer"`
			EdgeFollow       struct {
				Count int64 `json:"count"`
			} `json:"edge_follow"`
			FollowsViewer        bool   `json:"follows_viewer"`
			FullName             string `json:"full_name"`
			HasChannel           bool   `json:"has_channel"`
			HasBlockedViewer     bool   `json:"has_blocked_viewer"`
			HighlightReelCount   int    `json:"highlight_reel_count"`
			HasRequestedViewer   bool   `json:"has_requested_viewer"`
			ID                   string `json:"id"`
			IsBusinessAccount    bool   `json:"is_business_account"`
			IsJoinedRecently     bool   `json:"is_joined_recently"`
			BusinessCategoryName string `json:"business_category_name"`
			IsPrivate            bool   `json:"is_private"`
			IsVerified           bool   `json:"is_verified"`
			EdgeMutualFollowedBy struct {
				Count int64         `json:"count"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_mutual_followed_by"`
			ProfilePicURL          string      `json:"profile_pic_url"`
			ProfilePicURLHd        string      `json:"profile_pic_url_hd"`
			RequestedByViewer      bool        `json:"requested_by_viewer"`
			Username               string      `json:"username"`
			ConnectedFbPage        interface{} `json:"connected_fb_page"`
			EdgeFelixVideoTimeline struct {
				Count    int64         `json:"count"`
				PageInfo PageInfo      `json:"page_info"`
				Edges    []interface{} `json:"edges"`
			} `json:"edge_felix_video_timeline"`
			EdgeOwnerToTimelineMedia struct {
				Count    int64    `json:"count"`
				PageInfo PageInfo `json:"page_info"`
				Edges    []struct {
					Node struct {
						Typename           string `json:"__typename"`
						ID                 string `json:"id"`
						EdgeMediaToCaption struct {
							Edges []struct {
								Node struct {
									Text string `json:"text"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"edge_media_to_caption"`
						Shortcode          string `json:"shortcode"`
						EdgeMediaToComment struct {
							Count int64 `json:"count"`
						} `json:"edge_media_to_comment"`
						CommentsDisabled bool  `json:"comments_disabled"`
						TakenAtTimestamp int64 `json:"taken_at_timestamp"`
						Dimensions       struct {
							Height int `json:"height"`
							Width  int `json:"width"`
						} `json:"dimensions"`
						DisplayURL  string `json:"display_url"`
						EdgeLikedBy struct {
							Count int64 `json:"count"`
						} `json:"edge_liked_by"`
						EdgeMediaPreviewLike struct {
							Count int64 `json:"count"`
						} `json:"edge_media_preview_like"`
						Location struct {
							ID            string `json:"id"`
							HasPublicPage bool   `json:"has_public_page"`
							Name          string `json:"name"`
							Slug          string `json:"slug"`
						} `json:"location"`
						GatingInfo   interface{} `json:"gating_info"`
						MediaPreview string      `json:"media_preview"`
						Owner        struct {
							ID       string `json:"id"`
							Username string `json:"username"`
						} `json:"owner"`
						ThumbnailSrc       string `json:"thumbnail_src"`
						ThumbnailResources []struct {
							Src          string `json:"src"`
							ConfigWidth  int    `json:"config_width"`
							ConfigHeight int    `json:"config_height"`
						} `json:"thumbnail_resources"`
						IsVideo              bool   `json:"is_video"`
						AccessibilityCaption string `json:"accessibility_caption"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_owner_to_timeline_media"`
			EdgeSavedMedia struct {
				Count    int64         `json:"count"`
				PageInfo PageInfo      `json:"page_info"`
				Edges    []interface{} `json:"edges"`
			} `json:"edge_saved_media"`
			EdgeMediaCollections struct {
				Count    int64         `json:"count"`
				PageInfo PageInfo      `json:"page_info"`
				Edges    []interface{} `json:"edges"`
			} `json:"edge_media_collections"`
		} `json:"user"`
	} `json:"graphql"`
	FelixOnBoardingVideoResources struct {
		Mp4    string `json:"mp4"`
		Poster string `json:"poster"`
	} `json:"felix_onboarding_video_resources"`
	ShowFollowDialog bool `json:"show_follow_dialog"`
}

type ShortCodeMediaLikersResponse struct {
	Data struct {
		ShortCodeMedia *struct {
			ID          string `json:"id"`
			ShortCode   string `json:"shortcode"`
			EdgeLikedBy struct {
				Count    int64    `json:"count"`
				PageInfo PageInfo `json:"page_info"`
				Edges    []struct {
					Node struct {
						ID                string `json:"id"`
						Username          string `json:"username"`
						FullName          string `json:"full_name"`
						ProfilePicURL     string `json:"profile_pic_url"`
						IsPrivate         bool   `json:"is_private"`
						IsVerified        bool   `json:"is_verified"`
						FollowedByViewer  bool   `json:"followed_by_viewer"`
						RequestedByViewer bool   `json:"requested_by_viewer"`
						Reel              struct {
							ID              string      `json:"id"`
							ExpiringAt      int64       `json:"expiring_at"`
							LatestReelMedia interface{} `json:"latest_reel_media"`
							Seen            interface{} `json:"seen"`
							Owner           struct {
								Typename      string `json:"__typename"`
								ID            string `json:"id"`
								ProfilePicURL string `json:"profile_pic_url"`
								Username      string `json:"username"`
							} `json:"owner"`
						} `json:"reel"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_liked_by"`
		} `json:"shortcode_media"`
	} `json:"data"`
	Status string `json:"status"`
}

type WatchedStoryEntry struct {
	Item     StoryItem `json:"item" bson:"item"`
	Duration [2]int64  `json:"duration" bson:"duration"`
}

func MakeWatchedStoryRequest(belongs *Story, stories ...WatchedStoryEntry) map[string][]string {
	var m = make(map[string][]string)
	for _, story := range stories {
		m[fmt.Sprintf("%s_%s", story.Item.ID, belongs.ID)] = []string{fmt.Sprintf("%d_%d", story.Duration[0], story.Duration[1])}
	}
	return m
}

type HashTagStoriesResponse struct {
	Story  *Story `json:"story"`
	Status string `json:"status"`
}

type UserStoriesResponse struct {
	Broadcast interface{} `json:"broadcast"`
	Reel      Reel        `json:"reel"`
	Status    string      `json:"status"`
}
type FriendshipStatus struct {
	Following       bool `json:"following"`
	FollowedBy      bool `json:"followed_by"`
	Blocking        bool `json:"blocking"`
	Muting          bool `json:"muting"`
	IsPrivate       bool `json:"is_private"`
	IncomingRequest bool `json:"incoming_request"`
	OutgoingRequest bool `json:"outgoing_request"`
	IsBestie        bool `json:"is_bestie"`
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

func (s *StoryItem) GetOwner() *data.Owner {
	return &data.Owner{
		Id:       strconv.FormatInt(s.User.Pk, 10),
		Username: s.User.Username,
	}
}

type Reel struct {
	ID              int         `json:"id"`
	LatestReelMedia int         `json:"latest_reel_media"`
	ExpiringAt      int64       `json:"expiring_at"`
	Seen            int         `json:"seen"`
	CanReply        bool        `json:"can_reply"`
	CanReshare      bool        `json:"can_reshare"`
	ReelType        string      `json:"reel_type"`
	User            User        `json:"user"`
	Items           []StoryItem `json:"items"`
	PrefetchCount   int         `json:"prefetch_count"`
	HasBestiesMedia bool        `json:"has_besties_media"`
	MediaCount      int         `json:"media_count"`
}

type Owner struct {
	Type               string `json:"type"`
	Pk                 string `json:"pk"`
	Name               string `json:"name"`
	ProfilePicURL      string `json:"profile_pic_url"`
	ProfilePicUsername string `json:"profile_pic_username"`
}

type Story struct {
	ID                  string      `json:"id"`
	LatestReelMedia     int         `json:"latest_reel_media"`
	ExpiringAt          int64       `json:"expiring_at"`
	Seen                int         `json:"seen"`
	CanReply            bool        `json:"can_reply"`
	CanReshare          bool        `json:"can_reshare"`
	ReelType            string      `json:"reel_type"`
	Owner               Owner       `json:"owner"`
	Items               []StoryItem `json:"items"`
	PrefetchCount       int         `json:"prefetch_count"`
	UniqueIntegerReelID int64       `json:"unique_integer_reel_id"`
	Muted               bool        `json:"muted"`
}
