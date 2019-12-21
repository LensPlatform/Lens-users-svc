package models

type Settings struct {
	JsonEmbeddable
	Notifications Notification `json:"notifications" validate:"-" sql:"notification"`
	Privacy       Privacy      `json:"privacy" validate:"-" sql:"privacy"`
	Payment       Payments     `json:"payments" validate:"-" sql:"payments"`
}

type Payments struct {
	Debit       Card `json:"debit_card" validate:"-" sql:"debit_card"`
	Credit      Card `json:"credit_card" validate:"-" sql:"credit_card"`
	SecurityPin Pin  `json:"security_pin" validate:"-" sql:"security_pin"`
}

type Card struct {
	CardNumber   string `json:"card_number" validate:"required" sql:"card_number"`
	SecurityCode string `json:"security_code" validate:"required" sql:"security_code"`
	CardZipCode  string `json:"card_zip_code" validate:"required" sql:"card_zip_code"`
	FullName     string `json:"full_name" validate:"required" sql:"full_name"`
	Address      string `json:"address" validate:"required" sql:"address"`
	City         string `json:"city" validate:"required" sql:"city"`
	State        string `json:"state" validate:"required" sql:"state"`
	ZipCode      string `json:"zip_code" validate:"required" sql:"zip_code"`
}

type Pin struct {
	PinEnabled bool   `json:"pin_enabled" validate:"-" sql:"pin_enabled"`
	Pin        string `json:"pin" validate:"-" sql:"pin"`
}

type Privacy struct {
	ActivityStatus    bool     `json:"activity_status" validate:"-" sql:"activity_status"`
	PrivateAccount    bool     `json:"private" validate:"-" sql:"private"`
	BlockedAccountIds []string `json:"blocked_account_ids" validate:"-" sql:"blocked_account_ids"`
	MutedAccountIds   []string `json:"muted_account_ids" validate:"-" sql:"muted_account_ids"`
}

type Notification struct {
	JsonEmbeddable
	PauseAll              bool                                          `json:"pause_all" validate:"-" sql:"pause_all"`
	PostsAndComments      PostAndCommentsPushNotificationSettings       `json:"post_and_comments_pns" validate:"-" sql:"post_and_comments_pns"`
	FollowingAndFollowers FollowingAndFollowersPushNotificationSettings `json:"following_and_followers_pns" validate:"-" sql:"following_and_followers_pns"`
	DirectMessages        DirectMessagesPushNotificationSettings        `json:"direct_messages_pns" validate:"-" sql:"direct_messages_pns"`
	EmailAndSms           EmailAndSmsPushNotificationSettings           `json:"email_and_sms_pns" validate:"-" sql:"email_and_sms_pns"`
}

type PostAndCommentsPushNotificationSettings struct {
	JsonEmbeddable
	Likes                        TieredPushNotificationSetting `json:"likes_pns" validate:"-" sql:"likes_pns"`                                                           // yoan liked your photo
	LikesAndCommentsOnPostsOfYou TieredPushNotificationSetting `json:"likes_and_comments_on_posts_of_you_pns" validate:"-" sql:"likes_and_comments_on_posts_of_you_pns"` // yoan commented on a post you're tagged in
	PostsOfYou                   TieredPushNotificationSetting `json:"posts_of_you_pns" validate:"-" sql:"posts_of_you_pns"`                                             // yoan tagged you in a photo
	Comments                     TieredPushNotificationSetting `json:"comments_pns" validate:"-" sql:"comments_pns"`                                                     // yoan commented nice pic
	CommentLikes                 TieredPushNotificationSetting `json:"comment_likes_pns" validate:"-" sql:"comment_likes_pns"`                                           // yoan liked your comment "nice shot"
}

type FollowingAndFollowersPushNotificationSettings struct {
	JsonEmbeddable
	FollowerRequests         PushNotificationSetting       `json:"follower_requests_pns" validate:"-" sql:"follower_requests_pns"`                   // yoan has requested to follow you
	AcceptedFollowerRequests PushNotificationSetting       `json:"accepted_follower_requests_pns" validate:"-" sql:"accepted_follower_requests_pns"` // yoan accepted your follow request
	MentionsInBio            TieredPushNotificationSetting `json:"mentions_in_bio_pns" validate:"-" sql:"mentions_in_bio_pns"`                       // yoan mentioned you in his bio
}

type DirectMessagesPushNotificationSettings struct {
	JsonEmbeddable
	MessageRequests PushNotificationSetting `json:"message_requests_pns" validate:"-" sql:"message_requests_pns"`             // yoan wants to send you message
	Message         PushNotificationSetting `json:"message_pns" validate:"-" sql:"message_pns"`                               // yoan has sent you a message
	GroupRequests   PushNotificationSetting `json:"group_message_requests_pns" validate:"-" sql:"group_message_requests_pns"` // yoan wants to add willy to your group
}

type EmailAndSmsPushNotificationSettings struct {
	JsonEmbeddable
	FeedbackEmail  PushNotificationSetting `json:"feedback_emails_pns" validate:"-" sql:"feedback_emails_pns"`
	ReminderEmails PushNotificationSetting `json:"reminder_emails_pns" validate:"-" sql:"reminder_emails_pns"`
	ProductEmails  PushNotificationSetting `json:"product_emails_pns" validate:"-" sql:"product_emails_pns"`
	NewsEmails     PushNotificationSetting `json:"product_emails_pns" validate:"-" sql:"product_emails_pns"`
}

type TieredPushNotificationSetting struct {
	JsonEmbeddable
	Off               bool `json:"Off" validate:"-" sql:"Off"`
	FromPeopleIFollow bool `json:"from_people_l_follow" validate:"-" sql:"from_people_l_follow"`
	FromEveryone      bool `json:"from_everyone" validate:"-" sql:"from_everyone"`
}

type PushNotificationSetting struct {
	JsonEmbeddable
	Off bool `json:"Off" validate:"-" sql:"Off"`
	On  bool `json:"On" validate:"-" sql:"On"`
}
