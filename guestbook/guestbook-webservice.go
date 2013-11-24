// Copyright 2013 Bruno Albuquerque (bga@bug-br.org.br).
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

package guestbook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/codegangsta/martini"
)

func (g *GuestBook) GetPath() string {
	return "/guestbook"
}

func (g *GuestBook) WebDelete(params martini.Params) (int, string) {
	if len(params) == 0 {
		g.RemoveAllEntries()
		return 200, "collection deleted"
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return 400, "invalid entry id"
	}

	err = g.RemoveEntry(id)
	if err != nil {
		return 400, "entry not found"
	}

	return 200, "entry deleted"
}

func (g *GuestBook) WebGet(params martini.Params) (int, string) {
	if len(params) == 0 {
		encodedEntries, err := json.Marshal(g.GetAllEntries())
		if err != nil {
			return 500, "internal error"
		}
		return 200, string(encodedEntries)
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return 400, "invalid entry id"
	}

	entry, err := g.GetEntry(id)
	if err != nil {
		return 404, "entry not found"
	}

	encodedEntry, err := json.Marshal(entry)
	if err != nil {
		return 500, "internal error"
	}

	return 200, string(encodedEntry)
}

func (g *GuestBook) WebPost(params martini.Params,
	req *http.Request) (int, string) {
	defer req.Body.Close()
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return 500, "internal error"
	}

	if len(params) != 0 {
		return 405, "method not supported"
	}

	var guestBookEntry GuestBookEntry
	err = json.Unmarshal(requestBody, &guestBookEntry)
	if err != nil {
		return 400, "invalid JSON data"
	}

	g.AddEntry(guestBookEntry.Email, guestBookEntry.Title,
		guestBookEntry.Content)

	return 200, "new entry created"
}

