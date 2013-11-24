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

package webservice

import (
	"net/http"

	"github.com/codegangsta/martini"
)

type WebService interface {
	GetPath() string

	WebDelete(params martini.Params) (int, string)
	WebGet(params martini.Params) (int, string)
	WebPost(params martini.Params, req *http.Request) (int, string)
}

func RegisterWebService(webService WebService,
	classicMartini *martini.ClassicMartini) {
	path := webService.GetPath()

	classicMartini.Get(path, webService.WebGet)
	classicMartini.Get(path+"/:id", webService.WebGet)

	classicMartini.Post(path, webService.WebPost)
	classicMartini.Post(path+"/:id", webService.WebPost)

	classicMartini.Delete(path, webService.WebDelete)
	classicMartini.Delete(path+"/:id", webService.WebDelete)
}

