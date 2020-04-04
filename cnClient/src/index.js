// Copyright 2020 PerceptiSys Ltd, (Stephen Gaito)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import m from "mithril";

// Import JavaScript descriptions of the the interfaces
//
var action    = require("./interfaces/action")
var discovery = require("./interfaces/discovery")
var control   = require("./interfaces/control")
//
// Now import the main view
//
var Layout    = require("./views/Layout")

m.route(document.body, "/default", {
  ...action.routes,
  ...discovery.routes,
  ...control.routes,
  view: function() {
    return m(Layout)
  }
})
