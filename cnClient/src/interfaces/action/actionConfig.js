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

var ActionConfig = {
  data: {},
  storeData: function(result) {
    console.dir(result, {depth: null, colors: true});
    if (result == null) {
      result = {
      };
    }
    ActionConfig.data = result;
  },
  oninit: function(vnode) {
    m.request({
      method: "GET",
      url: "/action/".concat(vnode.attrs.actionName)
    }).then(function(result) {
      ActionConfig.storeData(result);
    })
  },
  argIterator: function() {
    console.dir(ActionConfig.data, { depth: null, colors: true});
    result = [];
    for (var key in ActionConfig.data.Args) {
      result.push(
        m("li", 
          m("strong", ActionConfig.data.Args[key].Key),
          m("p", ActionConfig.data.Args[key].Desc)
        )
      )
    }
    return result;
  },
  envIterator: function() {
    console.dir(ActionConfig.data, { depth: null, colors: true});
    result = [];
    for (var key in ActionConfig.data.Envs) {
      result.push(
        m("li", 
          m("strong", ActionConfig.data.Envs[key].Key),
          m("p", ActionConfig.data.Envs[key].Desc)
        )
      )
    }
    return result;
  },
  view: function(vnode) {
    return m("main.layout", [
      m("h1", ActionConfig.data.Name),
      m("p",  ActionConfig.data.Desc),
      m("h3", "Command line arguments:"),
      m("ul", ActionConfig.argIterator()),
      m("h3", "Environment variables:"),
      m("ul", ActionConfig.envIterator())
    ])
  },
  addRoutes: function(routes) {
    routes["/action/:actionName"] = {
      view: function(vnode) {
        return m(layout, m(ActionConfig))
      }
    }
  }
}

module.exports = ActionConfig