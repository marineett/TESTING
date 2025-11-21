package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func (s *APISuite) TestCreateUsersAndChats(t provider.T) {
	var (
		ctx            = context.Background()
		clientID       int64
		repetitorID    int64
		moderatorID    int64
		clientToken    string
		repetitorToken string
		cClient        *APIClient
		cRep           *APIClient
	)

	t.WithNewStep("Arrange", func(sx provider.StepCtx) {})

	t.WithNewStep("Act", func(sx provider.StepCtx) {
		type chatResp struct {
			ID int64 `json:"id"`
		}
		var authResp struct {
			Token  string `json:"token"`
			Role   string `json:"role"`
			UserID int64  `json:"user_id"`
		}

		resp, err := s.c.makeRequestWithBody(ctx, "POST", "/api/v2/registration", testClientData)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		sx.Require().NoError(json.Unmarshal(b, &authResp))
		clientID = authResp.UserID
		clientToken = authResp.Token

		resp, err = s.c.makeRequestWithBody(ctx, "POST", "/api/v2/registration", testRepetitorData)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)
		b, err = io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		sx.Require().NoError(json.Unmarshal(b, &authResp))
		repetitorID = authResp.UserID
		repetitorToken = authResp.Token

		resp, err = s.c.makeRequestWithBody(ctx, "POST", "/api/v2/registration", testModeratorData)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)
		b, err = io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		sx.Require().NoError(json.Unmarshal(b, &authResp))
		moderatorID = authResp.UserID

		cClient = s.c.WithToken(clientToken)
		cRep = s.c.WithToken(repetitorToken)

		body := map[string]interface{}{
			"type":         "client_repetitor",
			"client_id":    clientID,
			"repetitor_id": repetitorID,
			"moderator_id": 0,
		}
		resp, err = cClient.makeRequestWithBody(ctx, "POST", "/api/v2/chats", body)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)
		var crResp chatResp
		b, err = io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		sx.Require().NoError(json.Unmarshal(b, &crResp))
		crChatID := crResp.ID

		body = map[string]interface{}{
			"type":         "client_moderator",
			"client_id":    clientID,
			"moderator_id": moderatorID,
		}
		resp, err = cClient.makeRequestWithBody(ctx, "POST", "/api/v2/chats", body)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)
		var cmResp chatResp
		b, err = io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		sx.Require().NoError(json.Unmarshal(b, &cmResp))

		body = map[string]interface{}{
			"type":         "repetitor_moderator",
			"repetitor_id": repetitorID,
			"moderator_id": moderatorID,
		}
		resp, err = cRep.makeRequestWithBody(ctx, "POST", "/api/v2/chats", body)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)
		var rmResp chatResp
		b, err = io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		sx.Require().NoError(json.Unmarshal(b, &rmResp))

		resp, err = cClient.makeRequest(ctx, "PUT", fmt.Sprintf("/api/v2/chats/%d", crChatID), nil)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)

		msg := map[string]interface{}{
			"senderId": clientID,
			"content":  "test",
		}
		resp, err = cClient.makeRequestWithBody(ctx, "POST", fmt.Sprintf("/api/v2/chats/%d/messages", crChatID), msg)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)

		resp, err = cClient.makeRequest(ctx, "DELETE", fmt.Sprintf("/api/v2/chats/%d", crChatID), nil)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusNoContent, resp.StatusCode)
	})

	t.WithNewStep("Assert", func(sx provider.StepCtx) {})
}

func TestRunAPISuite(t *testing.T) {
	suite.RunSuite(t, new(APISuite))
}
