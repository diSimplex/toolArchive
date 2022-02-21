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

var layout = require("../../layout")

var Action = {
  data: {},
  storeData: function(result) {
    console.dir(result, {depth: null, colors: true});
    if (result == null) {
      result = {
      };
    }
    Action.data = result;
  },
  oninit: function() {
    m.request({
      method: "GET",
      url: "/action"
    }).then(function(result) {
      Action.storeData(result);
    })
  },
  dataIterator: function() {
    console.dir(Action.data, { depth: null, colors: true});
    result = [];
    for (var key in Action.data) {
      result.push(
        m("li", 
          m("strong",
            m(m.route.Link,
              { href: "/action/".concat(key) },
              Action.data[key].Name
            )
          ),
          m("p", Action.data[key].Desc)
        )
      )
    }
    return result;
  },
  view: function(vnode) {
    return m("main.layout", [
      m("h1", "Available  Actions"),
      m("ul", Action.dataIterator())
    ])
  },
  addRoutes: function(routes) {
    routes["/action"] = {
      view: function() {
        return m(layout, m(Action))
      }
    }
  }
}

module.exports = Action