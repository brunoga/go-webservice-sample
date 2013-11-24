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
	"fmt"
	"sync"
)

type GuestBookEntry struct {
	Id      int
	Email   string
	Title   string
	Content string
}

type GuestBook struct {
	guestBookData []*GuestBookEntry
	mutex *sync.Mutex
}

func NewGuestBook() *GuestBook {
	return &GuestBook{
		make([]*GuestBookEntry, 0),
		new(sync.Mutex),
	}
}

func (g *GuestBook) AddEntry(email, title, content string) int {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	newId := len(g.guestBookData)

	newEntry := &GuestBookEntry{
		newId,
		email,
		title,
		content,
	}

	g.guestBookData = append(g.guestBookData, newEntry)

	return newId
}

func (g *GuestBook) RemoveEntry(id int) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if id < 0 || id >= len(g.guestBookData) {
		return fmt.Errorf("invalid id")
	}

	g.guestBookData[id] = nil

	return nil
}

func (g *GuestBook) GetEntry(id int) (*GuestBookEntry, error) {
	if id < 0 || id >= len(g.guestBookData) ||
		g.guestBookData[id] == nil {
		return nil, fmt.Errorf("invalid id")
	}

	return g.guestBookData[id], nil
}

func (g *GuestBook) GetAllEntries() []*GuestBookEntry {
	entries := make([]*GuestBookEntry, 0)
	for _, entry := range g.guestBookData {
		if entry != nil {
			entries = append(entries, entry)
		}
	}

	return entries
}

func (g *GuestBook) RemoveAllEntries() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.guestBookData = []*GuestBookEntry{}
}

