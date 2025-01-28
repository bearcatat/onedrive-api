package resources

type DriveItem struct {
	DownloadURL string `json:"@microsoft.graph.downloadUrl,omitempty"`

	Audio       *Audio          `json:"audio,omitempty"`
	CTag        string          `json:"cTag,omitempty"`
	Deleted     *DeletedFacet   `json:"deleted,omitempty"`
	Description string          `json:"description,omitempty"`
	ETag        string          `json:"eTag,omitempty"`
	File        *File           `json:"file,omitempty"`
	Folder      *Folder         `json:"folder,omitempty"`
	Id          string          `json:"id,omitempty"`
	Image       *Image          `json:"image,omitempty"`
	Location    *GeoCoordinates `json:"location,omitempty"`
	Name        string          `json:"name,omitempty"`
	//Todo RemoteItem
	Photo  *Photo `json:"photo,omitempty"`
	Size   int64  `json:"size,omitempty"`
	Video  *Video `json:"video,omitempty"`
	WebURL string `json:"webUrl,omitempty"`
}

type Photo struct {
	CameraMake          string  `json:"cameraMake,omitempty"`
	CameraModel         string  `json:"cameraModel,omitempty"`
	TakenDateTime       string  `json:"takenDateTime,omitempty"`
	FNumber             float64 `json:"fNumber,omitempty"`
	ExposureDenominator float64 `json:"exposureDenominator,omitempty"`
	ExposureNumerator   float64 `json:"exposureNumerator,omitempty"`
	FocalLength         float64 `json:"focalLength,omitempty"`
	Iso                 int     `json:"iso,omitempty"`
}

// Video represents the video metadata of a OneDrive drive item.
type Video struct {
	Duration              int     `json:"duration,omitempty"`
	Height                float64 `json:"height,omitempty"`
	Width                 float64 `json:"width,omitempty"`
	AudioBitsPerSample    int     `json:"audioBitsPerSample,omitempty"`
	AudioChannels         int     `json:"audioChannels,omitempty"`
	AudioFormat           string  `json:"audioFormat,omitempty"`
	AudioSamplesPerSecond int     `json:"audioSamplesPerSecond,omitempty"`
	Bitrate               int     `json:"bitrate,omitempty"`
	FourCC                string  `json:"fourCC,omitempty"`
	FrameRate             float64 `json:"frameRate,omitempty"`
}

type Audio struct {
	Album             string `json:"album,omitempty"`
	AlbumArtist       string `json:"albumArtist,omitempty"`
	Artist            string `json:"artist,omitempty"`
	Bitrate           int    `json:"bitrate,omitempty"`
	Composers         string `json:"composers,omitempty"`
	Copyright         string `json:"copyright,omitempty"`
	Disc              int    `json:"disc,omitempty"`
	DiscCount         int    `json:"discCount,omitempty"`
	Duration          int    `json:"duration,omitempty"`
	Genre             string `json:"genre,omitempty"`
	HasDrm            bool   `json:"hasDrm,omitempty"`
	IsVariableBitrate bool   `json:"isVariableBitrate,omitempty"`
	Title             string `json:"title,omitempty"`
	Track             int    `json:"track,omitempty"`
	TrackCount        int    `json:"trackCount,omitempty"`
	Year              int    `json:"year,omitempty"`
}

type File struct {
	Hashes             Hashes `json:"hashes,omitempty"`
	MIMEType           string `json:"mimeType,omitempty"`
	ProcessingMetadata bool   `json:"processingMetadata,omitempty"`
}

// Folder
// childCount	Int32	Number of children contained immediately within this container.
// view	folderView	A collection of properties defining the recommended view for the folder.
type Folder struct {
	ChildCount int        `json:"childCount,omitempty"`
	View       FolderView `json:"view,omitempty"`
}

type Image struct {
	Height float64 `json:"height,omitempty"`
	Width  float64 `json:"width,omitempty"`
}

type DeletedFacet struct {
	State string `json:"state,omitempty"`
}

type GeoCoordinates struct {
	Altitude  float64 `json:"altitude,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type Hashes struct {
	Sha1Hash     string `json:"sha1Hash,omitempty"`
	Crc32Hash    string `json:"crc32Hash,omitempty"`
	QuickXorHash string `json:"quickXorHash,omitempty"`
}

type FolderView struct {
	SortBy    string `json:"sortBy,omitempty"`
	SortOrder string `json:"sortOrder,omitempty"`
	ViewType  string `json:"viewType,omitempty"`
}

func NewCreateFolderRequest(folderName string) *CreateFolderRequest {
	return &CreateFolderRequest{
		FolderName:       folderName,
		FolderFacet:      Facet{},
		ConflictBehavior: "rename",
	}
}

type CreateFolderRequest struct {
	FolderName       string `json:"name,omitempty"`
	FolderFacet      Facet  `json:"folder,omitempty"`
	ConflictBehavior string `json:"@microsoft.graph.conflictBehavior,omitempty"`
}

type Facet struct {
}

func NewUploadSessionRequest() *UploadSessionRequest {
	return &UploadSessionRequest{
		Item: UploadSessionRequestItem{
			ConflictBehavior: "rename",
		},
		DeferCommit: false,
	}
}

type UploadSessionRequest struct {
	Item        UploadSessionRequestItem `json:"item,omitempty"`
	DeferCommit bool                     `json:"deferCommit,omitempty"`
}

type UploadSessionRequestItem struct {
	// FileName         string `json:"name,omitempty"`
	ConflictBehavior string `json:"@microsoft.graph.conflictBehavior,omitempty"`
}

type UploadSession struct {
	UploadURL          string `json:"uploadUrl,omitempty"`
	ExpirationDateTime string `json:"expirationDateTime,omitempty"`
}

type UploadSessionResponse struct {
	ExpirationDateTime string   `json:"expirationDateTime,omitempty"`
	NextExpectedRanges []string `json:"nextExpectedRanges,omitempty"`
	DriveItem
}

type CopyRequest struct {
	ParentReference *ItemReference `json:"parentReference,omitempty"`
	Name            string         `json:"name,omitempty"`
}

func NewCopyRequest(parentItem *DriveItem, drive *Drive) *CopyRequest {
	return &CopyRequest{
		ParentReference: &ItemReference{
			ID:      parentItem.Id,
			DriveID: drive.Id,
		},
		Name: parentItem.Name,
	}
}

type ItemReference struct {
	DriveID   string `json:"driveId,omitempty"`
	DriveType string `json:"driveType,omitempty"`
	ID        string `json:"id,omitempty"`
	ListID    string `json:"listId,omitempty"`
	Name      string `json:"name,omitempty"`
	Path      string `json:"path,omitempty"`
	ShareID   string `json:"shareId,omitempty"`
	SiteID    string `json:"siteId,omitempty"`
}
