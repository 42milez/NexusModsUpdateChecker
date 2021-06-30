package error

const (
	CastFailed Err = iota
	CloseConnectionFailed
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
	DirectoryAlreadyExists
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
	ListArchivedContentFailed
	MoveFileFailed
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
	UnsupportedCommand
	UnsupportedVersionFormat
	UpdateRefFailed
)

var errors = map[Err]string{
	CastFailed:                "CAST_FAILED",
	CloseConnectionFailed:     "CLOSE_CONNECTION_FAILED",
	CloseFileFailed:           "CLOSE_FILE_FAILED",
	CopyFileFailed:            "COPY_FILE_FAILED",
	CreateCommitFailed:        "CREATE_COMMIT_FAILED",
	CreateDirectoryFailed:     "CREATE_DIRECTORY_FAILED",
	CreateFileFailed:          "CREATE_FILE_FAILED",
	CreatePullRequestFailed:   "CREATE_PULL_REQUEST_FAILED",
	CreateRefFailed:           "CREATE_REF_FAILED",
	CreateReleaseNoteFailed:   "CREATE_RELEASE_NOTE_FAILED",
	CreateRequestFailed:       "CREATE_REQUEST_FAILED",
	CreateTreeFailed:          "CREATE_TREE_FAILED",
	CtxCanceled:               "CTX_CANCELED_FAILED",
	DecodeJsonFailed:          "DECODE_JSON_FAILED",
	DecodeYamlFailed:          "DECODE_YAML_FAILED",
	DirectoryAlreadyExists:    "DIRECTORY_ALREADY_EXISTS",
	DownloadFailed:            "DOWNLOAD_FAILED",
	EncodeYamlFailed:          "ENCODE_YAML_FAILED",
	ExtractFailed:             "EXTRACT_FAILED",
	FileSizeMismatch:          "FILE_SIZE_MISMATCH",
	GetDownloadLinkFailed:     "GET_DOWNLOAD_LINK_FAILED",
	GetLatestFileFailed:       "GET_LATEST_FILE_FAILED",
	GetRefFailed:              "GET_REF_FAILED",
	GetSecretFailed:           "GET_SECRET_FAILED",
	GetWorkingDirectoryFailed: "GET_WORKING_DIRECTORY_FAILED",
	InvalidAccessToken:        "INVALID_ACCESS_TOKEN",
	InvalidApiKey:             "INVALID_API_KEY",
	InvalidBranchName:         "INVALID_BRANCH_NAME",
	ListArchivedContentFailed: "LIST_ARCHIVED_CONTENT_FAILED",
	MoveFileFailed:            "MOVE_FILE_FAILED",
	NoBranchProvided:          "NO_BRANCH_PROVIDED",
	NoCommitMessageProvided:   "NO_COMMIT_MESSAGE_PROVIDED",
	NoDescriptionProvided:     "NO_DESCRIPTION_PROVIDED",
	NoDownloadLinkReceived:    "NO_DOWNLOAD_LINK_RECEIVED",
	NoFileProvided:            "NO_FILE_PROVIDED",
	NoFileReceived:            "NO_FILE_RECEIVED",
	NoSubjectProvided:         "NO_SUBJECT_PROVIDED",
	OpenFileFailed:            "OPEN_FILE_FAILED",
	ParseTimeFailed:           "PARSE_TIME_FAILED",
	PushCommitFailed:          "PUSH_COMMIT_FAILED",
	ReadFileFailed:            "READ_FILE_FAILED",
	RecognizeVersionFailed:    "RECOGNIZE_VERSION_FAILED",
	RequestFailed:             "DO_REQUEST_FAILED",
	SecretNotFound:            "FIND_SECRET_FAILED",
	UnsupportedArchiveFormat:  "UNSUPPORTED_ARCHIVE_FORMAT",
	UnsupportedCommand:        "UNSUPPORTED_COMMAND",
	UnsupportedVersionFormat:  "UNSUPPORTED_VERSION_FORMAT",
	UpdateRefFailed:           "UPDATE_REF_FAILED",
}

type Err int

func (v Err) Error() string {
	return v.String()
}

func (v Err) String() string {
	return errors[v]
}
