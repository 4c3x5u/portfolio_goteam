package board

//import (
//	"bytes"
//	"encoding/json"
//	"net/http"
//	"net/http/httptest"
//	"server/assert"
//	"testing"
//)
//
//func TestPATCHHandler(t *testing.T) {
//	validator := &fakePATCHReqValidator{}
//	sut := NewPATCHHandler(patchValidator)
//
//	t.Run("InvalidBoardID", func(t *testing.T) {
//		username := ""
//		boardID := ""
//		reqBody, err := json.Marshal(ReqBody{Name: ""})
//		if err != nil {
//			t.Fatal(err)
//		}
//		req, err := http.NewRequest(
//			http.MethodPatch, "", bytes.NewReader(reqBody),
//		)
//		if err != nil {
//			t.Fatal(err)
//		}
//		req.RequestURI = "?boardID=" + boardID
//		w := httptest.NewRecorder()
//
//		validator.OutErr = false
//
//		sut.Handle(w, req, username)
//		res := w.Result()
//
//		if err = assert.Equal(
//			http.StatusBadRequest, res.StatusCode,
//		); err != nil {
//			t.Error(err)
//		}
//	})
//}
