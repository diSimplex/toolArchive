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

var m = require("mithril")

module.exports = {
  view: function(vnode) {
    return m("main.layout", [
      m("nav.nav-menu", [
        m(m.route.Link, {href: "/home",      class: "nav-menu-link"}, "ConTeXt Nursery"),
        m(m.route.Link, {href: "/action", class: "nav-menu-link"}, "Actions"),
        m(m.route.Link, {href: "/heartbeat",  class: "nav-menu-link"}, "Discovery"),
        m(m.route.Link, {href: "/control",  class: "nav-menu-link"}, "Control")
      ]),
      m("hr"),
      m("section", vnode.children),
    ])
  }
}
