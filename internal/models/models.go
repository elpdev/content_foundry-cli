package models

// Brand represents a brand from the list endpoint.
type Brand struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// BrandDetail is the full brand from the show endpoint.
type BrandDetail struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Slug             string `json:"slug"`
	Description      string `json:"description"`
	VoiceGuidelines  string `json:"voice_guidelines"`
	TargetAudience   string `json:"target_audience"`
	KeyInfo          string `json:"key_info"`
	ContactInfo      string `json:"contact_info"`
	MissionStatement string `json:"mission_statement"`
	Values           string `json:"values"`
	VisualIdentity   string `json:"visual_identity"`
	ContentPillars   string `json:"content_pillars"`
	Competitors      string `json:"competitors"`
	DosAndDonts      string `json:"dos_and_donts"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// Label represents a platform label.
type Label struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	ExternalID  *string `json:"external_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// Source represents a content source.
type Source struct {
	ID              int64            `json:"id"`
	Name            string           `json:"name"`
	Type            string           `json:"type"`
	Config          map[string]any   `json:"config"`
	PollingSchedule string           `json:"polling_schedule"`
	Active          bool             `json:"active"`
	Prompt          string           `json:"prompt"`
	LastFetchedAt   string           `json:"last_fetched_at"`
	CreatedAt       string           `json:"created_at"`
	UpdatedAt       string           `json:"updated_at"`
}

// ContentItem represents a content item.
type ContentItem struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	SourceID    int64  `json:"source_id"`
	SourceURL   string `json:"source_url"`
	Status      string `json:"status"`
	ContentHash string `json:"content_hash"`
	FetchedAt   string `json:"fetched_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ContentItemDraft is a minimal draft reference shown in content item detail.
type ContentItemDraft struct {
	ID           int64  `json:"id"`
	PlatformID   int64  `json:"platform_id"`
	Status       string `json:"status"`
	Version      int    `json:"version"`
	ScheduledFor string `json:"scheduled_for"`
	CreatedAt    string `json:"created_at"`
}

// Platform represents a publishing platform.
// LLMModel represents an AI model available for content generation.
type LLMModel struct {
	ID              int64  `json:"id"`
	ModelID         string `json:"model_id"`
	Name            string `json:"name"`
	Provider        string `json:"provider"`
	Family          string `json:"family"`
	ContextWindow   int64  `json:"context_window"`
	MaxOutputTokens int64  `json:"max_output_tokens"`
}

// ModelName returns a display string for the model, falling back to model_id.
func (m *LLMModel) DisplayName() string {
	if m.Name != "" {
		return m.Name
	}
	return m.ModelID
}

type Platform struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	Type           string         `json:"type"`
	Slug           string         `json:"slug"`
	Active         bool           `json:"active"`
	PromptTemplate string         `json:"prompt_template"`
	ModelID        any            `json:"model_id"`
	Model          *LLMModel      `json:"model,omitempty"`
	Settings       map[string]any `json:"settings"`
	CreatedAt      string         `json:"created_at"`
	UpdatedAt      string         `json:"updated_at"`
}

// ModelDisplayName returns a human-readable model name for the platform.
func (p *Platform) ModelDisplayName() string {
	if p.Model != nil {
		return p.Model.DisplayName()
	}
	return ""
}

// Draft represents a content draft.
type Draft struct {
	ID            int64  `json:"id"`
	ContentItemID int64  `json:"content_item_id"`
	PlatformID    int64  `json:"platform_id"`
	Content       string `json:"content"`
	Version       int    `json:"version"`
	Status        string `json:"status"`
	ReviewedByID  *int64 `json:"reviewed_by_id"`
	ReviewedAt    string `json:"reviewed_at"`
	ScheduledFor  string `json:"scheduled_for"`
	ScheduledByID *int64 `json:"scheduled_by_id"`
	AssignedToID  *int64 `json:"assigned_to_id"`
	AssignedAt    string `json:"assigned_at"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// DraftComment represents a comment on a draft.
type DraftComment struct {
	ID        int64  `json:"id"`
	DraftID   int64  `json:"draft_id"`
	UserID    int64  `json:"user_id"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

// RevisionRequest represents a request for revision on a draft.
type RevisionRequest struct {
	ID        int64  `json:"id"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"created_at"`
}

// Publication represents a published draft.
type Publication struct {
	ID           int64  `json:"id"`
	DraftID      int64  `json:"draft_id"`
	Status       string `json:"status"`
	ExternalID   string `json:"external_id"`
	URL          string `json:"url"`
	PublishedAt  string `json:"published_at"`
	ErrorMessage string `json:"error_message"`
	CreatedAt    string `json:"created_at"`
}

// Persona represents an AI persona.
type Persona struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	RoleTitle    string    `json:"role_title"`
	Description  string    `json:"description"`
	SystemPrompt string    `json:"system_prompt"`
	AvatarEmoji  string    `json:"avatar_emoji"`
	ModelID      any       `json:"model_id"`
	Model        *LLMModel `json:"model,omitempty"`
	Active       bool      `json:"active"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

// ModelDisplayName returns a human-readable model name for the persona.
func (p *Persona) ModelDisplayName() string {
	if p.Model != nil {
		return p.Model.DisplayName()
	}
	return ""
}

// Chat represents a chat session.
type Chat struct {
	ID          int64  `json:"id"`
	BrandID     int64  `json:"brand_id"`
	UserID      int64  `json:"user_id"`
	PersonaID   *int64 `json:"persona_id"`
	SubjectType string `json:"subject_type"`
	SubjectID   *int64 `json:"subject_id"`
	ModelID     any    `json:"model_id"`
	Title       string `json:"title"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Message represents a chat message.
type Message struct {
	ID           int64  `json:"id"`
	ChatID       int64  `json:"chat_id"`
	Role         string `json:"role"`
	Content      string `json:"content"`
	ModelID      any    `json:"model_id"`
	InputTokens  int    `json:"input_tokens"`
	OutputTokens int    `json:"output_tokens"`
	CreatedAt    string `json:"created_at"`
}

// MediaSession represents a media generation session.
type MediaSession struct {
	ID        int64  `json:"id"`
	BrandID   int64  `json:"brand_id"`
	UserID    int64  `json:"user_id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// MediaTurn is a polymorphic turn within a media session (image, video, or audio).
type MediaTurn struct {
	ID              int64  `json:"id"`
	Type            string `json:"type"`
	MediaSessionID  int64  `json:"media_session_id"`
	Position        int    `json:"position"`
	UserPrompt      string `json:"user_prompt"`
	Status          string `json:"status"`
	AspectRatio     string `json:"aspect_ratio"`
	DurationSeconds int    `json:"duration_seconds"`
	ImageURL        string `json:"image_url"`
	VideoURL        string `json:"video_url"`
	AudioURL        string `json:"audio_url"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// ImageTurn represents an image generation turn.
type ImageTurn struct {
	ID             int64  `json:"id"`
	MediaSessionID int64  `json:"media_session_id"`
	Position       int    `json:"position"`
	UserPrompt     string `json:"user_prompt"`
	Status         string `json:"status"`
	AspectRatio    string `json:"aspect_ratio"`
	ImageURL       string `json:"image_url"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// VideoTurn represents a video generation turn.
type VideoTurn struct {
	ID              int64  `json:"id"`
	MediaSessionID  int64  `json:"media_session_id"`
	Position        int    `json:"position"`
	UserPrompt      string `json:"user_prompt"`
	Status          string `json:"status"`
	AspectRatio     string `json:"aspect_ratio"`
	DurationSeconds int    `json:"duration_seconds"`
	VideoURL        string `json:"video_url"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// AudioTurn represents an audio generation turn.
type AudioTurn struct {
	ID             int64   `json:"id"`
	MediaSessionID int64   `json:"media_session_id"`
	Position       int     `json:"position"`
	UserPrompt     string  `json:"user_prompt"`
	Status         string  `json:"status"`
	CaptureSeconds int     `json:"capture_seconds"`
	BPM            *int    `json:"bpm"`
	Density        *float64 `json:"density"`
	Brightness     *float64 `json:"brightness"`
	Scale          string  `json:"scale"`
	AudioURL       string  `json:"audio_url"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// Asset represents an uploaded media asset.
type Asset struct {
	ID          int64  `json:"id"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	ByteSize    int64  `json:"byte_size"`
	FileURL     string `json:"file_url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// BrandDocument represents a brand knowledge document.
type BrandDocument struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	SourceType   string `json:"source_type"`
	SourceURL    string `json:"source_url"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// ActivityEvent represents an activity log entry.
type ActivityEvent struct {
	ID            int64          `json:"id"`
	EventType     string         `json:"event_type"`
	TrackableType string         `json:"trackable_type"`
	TrackableID   int64          `json:"trackable_id"`
	UserID        int64          `json:"user_id"`
	Metadata      map[string]any `json:"metadata"`
	CreatedAt     string         `json:"created_at"`
}

// Notification represents a user notification.
type Notification struct {
	ID        int64          `json:"id"`
	Type      string         `json:"type"`
	Params    map[string]any `json:"params"`
	ReadAt    string         `json:"read_at"`
	CreatedAt string         `json:"created_at"`
}

// Ticket represents a support ticket.
type Ticket struct {
	ID             int64          `json:"id"`
	Subject        string         `json:"subject"`
	Description    string         `json:"description"`
	Category       string         `json:"category"`
	Status         string         `json:"status"`
	Priority       string         `json:"priority"`
	Metadata       map[string]any `json:"metadata"`
	AdminNotes     string         `json:"admin_notes"`
	ResolvedAt     string         `json:"resolved_at"`
	ResolvedReason string         `json:"resolved_reason"`
	CreatedAt      string         `json:"created_at"`
	UpdatedAt      string         `json:"updated_at"`
}

// Account represents the user's account.
type Account struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// AccountUser represents a team member.
type AccountUser struct {
	ID        int64    `json:"id"`
	UserID    int64    `json:"user_id"`
	AccountID int64    `json:"account_id"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// AccountInvitation represents a pending invitation.
type AccountInvitation struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Token       string   `json:"token"`
	Roles       []string `json:"roles"`
	BrandIDs    []int64  `json:"brand_ids"`
	InvitedByID int64    `json:"invited_by_id"`
	CreatedAt   string   `json:"created_at"`
}
