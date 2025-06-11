package api_tests

import (
	"testing"

	"pnxlr.eu.org/roll/net/up/api"
)

func TestRobotUploader(t *testing.T) {
	uploader := api.NewRobotUploader()

	raw := `{"statusCode":0,"msg":"操作成功","data":{"originText":"deadbabedeadbabedeadbabedeadbabe","resultText":null,"uploadType":"png","timestamp":1754321012345}}`
	data, err := uploader.Json([]byte(raw))
	if err != nil {
		t.Errorf("Json() error: %v", err)
	}
	if !uploader.Success(data) {
		t.Errorf("Success() error: expected %v, got %v", false, true)
	}
	if eid, id := "deadbabedeadbabedeadbabedeadbabe", uploader.ObjectID(data); id != eid {
		t.Errorf("ObjectID() error: expected %v, got %v", eid, id)
	}

	raw = `{"statusCode":0,"msg":"操作成功","data":{"originText":"不符合上传类型","resultText":null,"uploadType":"error","timestamp":0}}`
	data, err = uploader.Json([]byte(raw))
	if err != nil {
		t.Errorf("Json() error: %v", err)
	}
	if uploader.Success(data) {
		t.Errorf("Success() error: expected %v, got %v", false, true)
	}
}
