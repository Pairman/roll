package api

import (
	"encoding/json"
	"net/http"

	netUtil "pnxlr.eu.org/roll/net/util"
)

// RobotUploader stolen from "https://48sydjll.mh.chaoxing.com".
//
// Accepts: PNG, JPG, JPEG, DOC, DOCX, PDF.
//
// Limits: Up to 2GiB. Validates file extension and content.
type RobotUploader struct {
	URL     string
	Headers http.Header
}

type RobotUploadJson struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
	Data       struct {
		OriginText string  `json:"originText"`
		ResultText *string `json:"resultText"`
		UploadType string  `json:"uploadType"`
		Timestamp  int64   `json:"timestamp"`
	} `json:"data"`
}

func NewRobotUploader() *RobotUploader {
	return &RobotUploader{
		URL:     "https://robot.chaoxing.com/v1/front/uploadKnowledgeFile",
		Headers: netUtil.GlobalHeader.Clone(),
	}
}

func (uploader *RobotUploader) Json(raw []byte) (*RobotUploadJson, error) {
	data := &RobotUploadJson{}
	return data, json.Unmarshal(raw, data)
}

func (uploader *RobotUploader) Success(data *RobotUploadJson) bool {
	return data.StatusCode == 0 && data.Data.Timestamp != 0
}

func (uploader *RobotUploader) ObjectID(data *RobotUploadJson) string {
	return data.Data.OriginText
}
