package error

const (
	CloseConnectionFailed Err = iota
	CloseFileFailed
	CopyFileFailed
	CreateCommitFailed
	CreateDirectoryFailed
	CreateFileFailed
	CreatePullRequestFailed
	CreateRefFailed
	CreateReleaseNoteFailed
	CreateRequestFailed
	CreateTreeFailed
	CtxCanceled
	DecodeJsonFailed
	DecodeYamlFailed
	DownloadFailed
	EncodeYamlFailed
	ExtractFailed
	FileSizeMismatch
	GetDownloadLinkFailed
	GetLatestFileFailed
	GetRefFailed
	GetSecretFailed
	GetWorkingDirectoryFailed
	InvalidAccessToken
	InvalidApiKey
	InvalidBranchName
	NoBranchProvided
	NoCommitMessageProvided
	NoDescriptionProvided
	NoDownloadLinkReceived
	NoFileProvided
	NoFileReceived
	NoSubjectProvided
	OpenFileFailed
	ParseTimeFailed
	PushCommitFailed
	ReadFileFailed
	RecognizeVersionFailed
	RequestFailed
	SecretNotFound
	UnsupportedArchiveFormat
	UnsupportedVersionFormat
	UpdateRefFailed
)

var errors = map[Err]string{
	CloseConnectionFailed:     "CANT_CLOSE_CONNECTION",
	CloseFileFailed:           "CANT_CLOSE_FILE",
	CopyFileFailed:            "CANT_COPY_FILE",
	CreateCommitFailed:        "CANT_CREATE_COMMIT",
	CreateDirectoryFailed:     "CANT_CREATE_DIRECTORY",
	CreateFileFailed:          "CANT_CREATE_FILE",
	CreatePullRequestFailed:   "CANT_CREATE_PULL_REQUEST",
	CreateRefFailed:           "CANT_CREATE_REF",
	CreateReleaseNoteFailed:   "CANT_CREATE_RELEASE_NOTE",
	CreateRequestFailed:       "CANT_CREATE_REQUEST",
	CreateTreeFailed:          "CANT_CREATE_TREE",
	CtxCanceled:               "CTX_CANCELED",
	DecodeJsonFailed:          "CANT_DECODE_JSON",
	DecodeYamlFailed:          "CANT_DECODE_YAML",
	DownloadFailed:            "CANT_DOWNLOAD",
	EncodeYamlFailed:          "CANT_ENCODE_YAML",
	ExtractFailed:             "CANT_EXTRACT",
	FileSizeMismatch:          "FILE_SIZE_MISMATCH",
	GetDownloadLinkFailed:     "CANT_GET_DOWNLOAD_LINK",
	GetLatestFileFailed:       "CANT_GET_LATEST_FILE",
	GetRefFailed:              "CANT_GET_REF",
	GetSecretFailed:           "CANT_GET_SECRET",
	GetWorkingDirectoryFailed: "CANT_GET_WORKING_DIRECTORY",
	InvalidAccessToken:        "INVALID_ACCESS_TOKEN",
	InvalidApiKey:             "INVALID_API_KEY",
	InvalidBranchName:         "INVALID_BRANCH_NAME",
	NoBranchProvided:          "NO_BRANCH_PROVIDED",
	NoCommitMessageProvided:   "NO_COMMIT_MESSAGE_PROVIDED",
	NoDescriptionProvided:     "NO_DESCRIPTION_PROVIDED",
	NoDownloadLinkReceived:    "NO_DOWNLOAD_LINK_RECEIVED",
	NoFileProvided:            "NO_FILE_PROVIDED",
	NoFileReceived:            "NO_FILE_RECEIVED",
	NoSubjectProvided:         "NO_SUBJECT_PROVIDED",
	OpenFileFailed:            "CANT_OPEN_FILE",
	ParseTimeFailed:           "CANT_PARSE_TIME",
	PushCommitFailed:          "CANT_PUSH_COMMIT",
	ReadFileFailed:            "CANT_READ_FILE",
	RecognizeVersionFailed:    "CANT_RECOGNIZE_VERSION",
	RequestFailed:             "CANT_DO_REQUEST",
	SecretNotFound:            "CANT_FIND_SECRET",
	UnsupportedArchiveFormat:  "UNSUPPORTED_ARCHIVE_FORMAT",
	UnsupportedVersionFormat:  "UNSUPPORTED_VERSION_FORMAT",
	UpdateRefFailed:           "CANT_UPDATE_REF",
}

type Err int

func (v Err) Error() string {
	return v.String()
}

func (v Err) String() string {
	return errors[v]
}
