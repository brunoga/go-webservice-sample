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

// GuestBookEntry represents a single entry in a Guest Book. It contains the
// usual fields.
type GuestBookEntry struct {
	Id      int
	Email   string
	Title   string
	Content string
}

// GuestBook represents a Guest Book instance. It holds the associated
// GuestBookEntries.
type GuestBook struct {
	guestBookData []*GuestBookEntry
	mutex *sync.Mutex
}

// NewGuestBook returns a new empty GuestBook instance.
func NewGuestBook() *GuestBook {
	return &GuestBook{
		make([]*GuestBookEntry, 0),
		new(sync.Mutex),
	}
}

// AddEntry adds a new GuestBookEntry with the provided data.
func (g *GuestBook) AddEntry(email, title, content string) int {
	// Acquire our lock and make sure it will be released.
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Get an id for this entry.
	newId := len(g.guestBookData)

	// Create new entry with the given data and the computed newId.
	newEntry := &GuestBookEntry{
		newId,
		email,
		title,
		content,
	}

	// Add entry to the Guest Book.
	g.guestBookData = append(g.guestBookData, newEntry)

	// Return the Id for the new entry.
	return newId
}

// RemoveEntry removes the entry with the given id. Return nil in case of
// success or a specific error in case of failure.
func (g *GuestBook) RemoveEntry(id int) error {
	// Acquaire our lock and make sure it will be released.
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Check if we have a valid id.
	if id < 0 || id >= len(g.guestBookData) ||
		g.guestBookData[id] == nil {
		return fmt.Errorf("invalid id")
	}

	// Set entry to nil. This is not memory efficient, but it is simple to
	// do.
	g.guestBookData[id] = nil

	return nil
}

// GetEntry returns the entry identified by the given id or an error if it can
// not find it.
func (g *GuestBook) GetEntry(id int) (*GuestBookEntry, error) {
	// Check if we have a valid id.
	if id < 0 || id >= len(g.guestBookData) ||
		g.guestBookData[id] == nil {
		return nil, fmt.Errorf("invalid id")
	}

	// Return the associated entry.
	return g.guestBookData[id], nil
}

// GetAllEntries returns all non-nil entries in the Guest Book.
func (g *GuestBook) GetAllEntries() []*GuestBookEntry {
	// Placeholder for the entries we will be returning.
	entries := make([]*GuestBookEntry, 0)

	// Iterate through all existig entries.
	for _, entry := range g.guestBookData {
		if entry != nil {
			// Entry is not nil, so we want to return it.
			entries = append(entries, entry)
		}
	}

	return entries
}

// RemoveAllEntries removes all entries from the Guest Book.
func (g *GuestBook) RemoveAllEntries() {
	// Acquire our lock and make sure it will be released.
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Reset guestbook to a new empty one.
	g.guestBookData = []*GuestBookEntry{}
}

