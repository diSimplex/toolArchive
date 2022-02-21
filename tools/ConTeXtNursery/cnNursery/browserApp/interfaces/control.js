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

var layout = require("../layout")

var Control = {
  data: {},
  storeData: function(result) {
    console.dir(result, {depth: null, colors: true});
    if ( result == null ) {
      result = {
      };
    }
    Control.data = result;
  },
  oninit: function() {
    m.request({
      method: "GET",
      url: "/control"
    }).then(function(result) {
      Control.storeData(result);
    })
  },
  changeState: function(baseUrl, urlModifier, newState) {
    m.request({
      method: "PUT",
      url: baseUrl
        .concat("/control")
        .concat(urlModifier)
        .concat("/")
        .concat(newState)
    }).then(function(result) {
      Control.storeData(result)
    })
  },
  changeStateButton: function(key, newState, label) {
    return m("button", {
      onclick: function() {
        Control.changeState(
          Control.data[key].Base_Url,
          Control.data[key].Url_Modifier,
          newState
        )
      }
    }, label)
  },
  dataIterator: function() {
    console.dir(Control.data, { depth: null, colors: true });
    result = [
      m("tr", 
        m("th", { colspan: 2}),
        m("th", { colspan: 4}, "State")
      ),
      m("tr", 
        m("th", { colspan: 2}),
        m("th", { colspan: 4}, 
          m("hr", { style: { "margin-top": "0em", "margin-bottom": "0em" } } )
        )
      ),
      m("tr", 
        m("th", "Name"),
        m("th", "Processes"),
        m("th", "Current"),
        m("th", "Up"),
        m("th", "Pause"),
        m("th", "Kill")
      )
    ];
    for (var key in Control.data) {
      result.push(
        m("tr",
          m("td", key),
          m("td", Control.data[key].Processes),
          m("td", Control.data[key].State),
          m("td", Control.changeStateButton(key, "up",     "Up")),
          m("td", Control.changeStateButton(key, "paused", "Pause")),
          m("td", Control.changeStateButton(key, "kill",   "Kill"))
        )
      )
    }
    return result;
  },
  view: function(vnode) {
    return m("main.layout", [
      m("h1", "Federation Control Information"),
      m("table", Control.dataIterator())
    ])
  },
  addRoutes: function(routes) {
    routes["/control"] = {
      view: function() {
        return m(layout, m(Control))
      }
    }
  }
}

module.exports = Control