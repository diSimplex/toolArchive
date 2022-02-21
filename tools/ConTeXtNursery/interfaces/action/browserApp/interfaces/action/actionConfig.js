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
  actionName: "unknown",
  form: {
    Args: {},
    Envs: {}
  },
  runAction: function() {
    if (ActionConfig.actionName == "unknown") { return }
    
    console.dir(ActionConfig.form, {depth: null, colors: true})

    var actionConfig = { Args: [], Envs: [] }
    for (var key in ActionConfig.form.Args) {
      actionConfig.Args.push({
        Key: key,
        Value: ActionConfig.form.Args[key]
      })
    }
    for (var key in ActionConfig.form.Envs) {
      actionConfig.Envs.push({
        Key: key,
        Value: ActionConfig.form.Envs[key]
      })
    }
    console.dir(actionConfig, {depth: null, colors: true})

    m.request({
      method: "POST",
      url: "/action/".concat(ActionConfig.actionName),
      body: actionConfig
    })
  },
  data: {
    Name: "unknown",
    Desc: "unknown",
    Args: {},
    Envs: {}
  },
  oninit: function(vnode) {
    m.request({
      method: "GET",
      url: "/action/".concat(ActionConfig.actionName)
    }).then(function(result) {
    console.dir(result, {depth: null, colors: true});
    if (result == null) {
      result = {
      };
    }
    ActionConfig.data = result;
    })
  },
  argIterator: function() {
    console.dir(ActionConfig.data, { depth: null, colors: true});
    if (Object.keys(ActionConfig.data.Args).length < 1) { return null; }

    result = [];
    for (var key in ActionConfig.data.Args) {
      result.push(
        m("li", 
          m("strong", ActionConfig.data.Args[key].Key),
          " ",
          m("input", {
            type: "text",
            name: ActionConfig.data.Args[key].Key,
            onchange: function(evnt) {
              theKey = evnt.target.name
              ActionConfig.form.Args[theKey] = evnt.target.value
              //console.dir(evnt.target, {depth: null, colors: true})
              console.dir(ActionConfig.form, {depth: null, colors: true})
            }
          }),
          m("div", ActionConfig.data.Args[key].Desc)
        )
      )
    }
    return [
      m("h3", "Command line arguments:"),
      m("ul", result)
    ];
  },
  envIterator: function() {
    console.dir(ActionConfig.data, { depth: null, colors: true});
    if (Object.keys(ActionConfig.data.Envs).length < 1) { return null; }
    
    result = [];
    for (var key in ActionConfig.data.Envs) {
      result.push(
        m("li", 
          m("strong", ActionConfig.data.Envs[key].Key),
          " ",
          m("input", {
            type: "text",
            name: ActionConfig.data.Envs[key].Key,
            onchange: function(evnt) {
              theKey = evnt.target.name
              ActionConfig.form.Envs[theKey] = evnt.target.value;
              //console.dir(evnt.target, {depth: null, colors: true})
              console.dir(ActionConfig.form, {depth: null, colors: true})
            }
          }),
          m("div", ActionConfig.data.Envs[key].Desc)
        )
      )
    }
    return [
      m("h3", "Environment variables:"),
      m("ul", result)
    ];
  },
  view: function(vnode) {
    if (ActionConfig.data.Name == "unknown") { return m("main.layout") }
    
    return m("div", [
      m("h1", ActionConfig.data.Name),
      m("p",  ActionConfig.data.Desc),
      ActionConfig.argIterator(),
      ActionConfig.envIterator(),
      m("button", {
        onclick: function() {
          ActionConfig.runAction()
        }
      }, "Run")
    ])
  },
  addRoutes: function(routes) {
    routes["/action/:actionName"] = {
      view: function(vnode) {
        ActionConfig.actionName = vnode.attrs.actionName;
        return m(layout, m(ActionConfig))
      }
    }
  }
}

module.exports = ActionConfig