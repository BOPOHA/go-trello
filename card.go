/*
Copyright 2014 go-trello authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package trello

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Card struct {
	client                *Client
	Id                    string   `json:"id"`
	Name                  string   `json:"name"`
	Email                 string   `json:"email"`
	IdShort               int      `json:"idShort"`
	IdAttachmentCover     string   `json:"idAttachmentCover"`
	IdCheckLists          []string `json:"idCheckLists"`
	IdBoard               string   `json:"idBoard"`
	IdList                string   `json:"idList"`
	IdMembers             []string `json:"idMembers"`
	IdMembersVoted        []string `json:"idMembersVoted"`
	ManualCoverAttachment bool     `json:"manualCoverAttachment"`
	Closed                bool     `json:"closed"`
	Pos                   float32  `json:"pos"`
	ShortLink             string   `json:"shortLink"`
	DateLastActivity      string   `json:"dateLastActivity"`
	ShortUrl              string   `json:"shortUrl"`
	Subscribed            bool     `json:"subscribed"`
	Url                   string   `json:"url"`
	Due                   string   `json:"due"`
	Desc                  string   `json:"desc"`
	DescData              struct {
		Emoji struct{} `json:"emoji"`
	} `json:"descData"`
	CheckItemStates []struct {
		IdCheckItem string `json:"idCheckItem"`
		State       string `json:"state"`
	} `json:"checkItemStates"`
	Badges struct {
		Votes              int    `json:"votes"`
		ViewingMemberVoted bool   `json:"viewingMemberVoted"`
		Subscribed         bool   `json:"subscribed"`
		Fogbugz            string `json:"fogbugz"`
		CheckItems         int    `json:"checkItems"`
		CheckItemsChecked  int    `json:"checkItemsChecked"`
		Comments           int    `json:"comments"`
		Attachments        int    `json:"attachments"`
		Description        bool   `json:"description"`
		Due                string `json:"due"`
	} `json:"badges"`
	Labels []struct {
		Color string `json:"color"`
		Name  string `json:"name"`
	} `json:"labels"`
}

func (c *Client) Card(CardId string) (card *Card, err error) {
	body, err := c.Get("/card/" + CardId)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &card)
	card.client = c
	return
}

func (c *Card) Checklists() (checklists []Checklist, err error) {
	body, err := c.client.Get("/card/" + c.Id + "/checklists")
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &checklists)
	for i := range checklists {
		list := &checklists[i]
		list.client = c.client
		for i := range list.CheckItems {
			item := &list.CheckItems[i]
			item.client = c.client
			item.listID = list.Id
		}
	}
	return
}

func (c *Card) Members() (members []Member, err error) {
	body, err := c.client.Get("/cards/" + c.Id + "/members")
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &members)
	for i := range members {
		members[i].client = c.client
	}
	return
}

func (c *Card) Attachments() (attachments []Attachment, err error) {
	body, err := c.client.Get("/cards/" + c.Id + "/attachments")
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &attachments)
	for i := range attachments {
		attachments[i].client = c.client
	}
	return
}

// Attachment will return the specified attachment on the card
// https://developers.trello.com/advanced-reference/card#get-1-cards-card-id-or-shortlink-attachments-idattachment
func (c *Card) Attachment(attachmentId string) (*Attachment, error) {
	body, err := c.client.Get("/cards/" + c.Id + "/attachments/" + attachmentId)
	if err != nil {
		return nil, err
	}

	attachment := &Attachment{}
	err = json.Unmarshal(body, attachment)
	attachment.client = c.client
	return attachment, err
}

func (c *Card) Actions() (actions []Action, err error) {
	body, err := c.client.Get("/cards/" + c.Id + "/actions?filter=commentCard%2CcreateCard")
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &actions)
	for i := range actions {
		actions[i].client = c.client
	}
	return
}

// AddChecklist will add a checklist to the card.
// https://developers.trello.com/advanced-reference/card#post-1-cards-card-id-or-shortlink-checklists
func (c *Card) AddChecklist(name string) (*Checklist, error) {
	payload := url.Values{}
	payload.Set("name", name)
	body, err := c.client.Post("/cards/"+c.Id+"/checklists", payload)
	if err != nil {
		return nil, err
	}

	newList := &Checklist{}
	if err = json.Unmarshal(body, newList); err != nil {
		return nil, err
	}
	newList.client = c.client
	// the new list has no items, no need to walk those adding client
	return newList, err
}

// AddComment will add a new comment to the card
// https://developers.trello.com/advanced-reference/card#post-1-cards-card-id-or-shortlink-actions-comments
func (c *Card) AddComment(text string) (*Action, error) {
	payload := url.Values{}
	payload.Set("text", text)

	body, err := c.client.Post("/cards/"+c.Id+"/actions/comments", payload)
	if err != nil {
		return nil, err
	}

	newAction := &Action{}
	if err = json.Unmarshal(body, newAction); err != nil {
		return nil, err
	}
	newAction.client = c.client
	return newAction, nil
}

// Archive will archive the card
// https://developers.trello.com/advanced-reference/card#put-1-cards-card-id-or-shortlink-closed
func (c *Card) Archive() (*Card, error) {
	payload := url.Values{}
	payload.Set("value", "true")

	body, err := c.client.Put("/cards/"+c.Id+"/closed", payload)
	if err != nil {
		return nil, err
	}

	newCard := &Card{}
	if err = json.Unmarshal(body, newCard); err != nil {
		return nil, err
	}
	newCard.client = c.client
	return newCard, nil
}

// SendToBoard will dearchive the card, or send the card to the board back from archive
// https://developers.trello.com/advanced-reference/card#put-1-cards-card-id-or-shortlink-closed
func (c *Card) SendToBoard() (*Card, error) {
	payload := url.Values{}
	payload.Set("value", "false")

	body, err := c.client.Put("/cards/"+c.Id+"/closed", payload)
	if err != nil {
		return nil, err
	}
	newCard := &Card{}
	if err = json.Unmarshal(body, newCard); err != nil {
		return nil, err
	}
	newCard.client = c.client
	return newCard, nil
}

// MoveToList will move the card to another list
// https://developers.trello.com/advanced-reference/card#put-1-cards-card-id-or-shortlink-idlist
func (c *Card) MoveToList(listId string) (*Card, error) {
	payload := url.Values{}
	payload.Set("value", listId)

	body, err := c.client.Put("/cards/"+c.Id+"/idList", payload)
	if err != nil {
		return nil, err
	}
	newCard := &Card{}
	if err = json.Unmarshal(body, newCard); err != nil {
		return nil, err
	}
	newCard.client = c.client
	return newCard, nil
}

// MoveToPos will move card to the specified position
// https://developers.trello.com/advanced-reference/card#put-1-cards-card-id-or-shortlink-pos
func (c *Card) MoveToPos(pos int) (*Card, error) {
	payload := url.Values{}
	payload.Set("value", strconv.Itoa(pos))

	body, err := c.client.Put("/cards/"+c.Id+"/pos", payload)
	if err != nil {
		return nil, err
	}
	newCard := &Card{}
	if err = json.Unmarshal(body, newCard); err != nil {
		return nil, err
	}
	newCard.client = c.client
	return newCard, nil
}

// AddLableByID will add a label to a card
// https://developers.trello.com/reference#label-object
func (c *Card) AddLableByID(idLabel string) (*[]string, error) {
	payload := url.Values{}
	payload.Set("value", idLabel)

	body, err := c.client.Post("/cards/"+c.Id+"/idLabels", payload)
	if err != nil {
		return nil, err
	}

	newIDs := &[]string{}
	if err = json.Unmarshal(body, newIDs); err != nil {
		return nil, err
	}

	return newIDs, nil
}
