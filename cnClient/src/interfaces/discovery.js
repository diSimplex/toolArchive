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

var Discovery = {
  data: {},
  oninit: function() {
    m.request({
      method: "GET",
      url: "/heartbeat",
    }).then(function(result) {
      console.dir(result, {depth: null, colors: true});
      Discovery.data = result;
    })
  },
  dataIterator: function() {
    result = [];
    header = m("tr",
      m("th", "Name"),
      m("th", "Port"),
      m("th", "State"),
      m("th", "Processes"),
      m("th", "Cores"),
      m("th", "Speed Mhz"),
      m("th", "Mem Total"),
      m("th", "Mem Used"),
      m("th", "Swap Total"),
      m("th", "Swap Used"),
      m("th", "Load 1 min"),
      m("th", "Load 5 min"),
      m("th", "Load 15 min"),
    );
    result.push(header);

    for (var key in Discovery.data) {
      result.push(
        m("tr",
          m("td", 
            m("a.link",
              { href: Discovery.data[key].Base_Url },
              Discovery.data[key].Name                
            )
          ),
          m("td", Discovery.data[key].Port),
          m("td", Discovery.data[key].State),
          m("td", Discovery.data[key].Processes),
          m("td", Discovery.data[key].Cores),
          m("td", Discovery.data[key].Speed_Mhz),
          m("td", Discovery.data[key].Memory.Total),
          m("td", Discovery.data[key].Memory.Used),
          m("td", Discovery.data[key].Swap.Total),
          m("td", Discovery.data[key].Swap.Used),
          m("td", Discovery.data[key].Load.Load1),
          m("td", Discovery.data[key].Load.Load5),
          m("td", Discovery.data[key].Load.Load15),            
        )
      );
    }
    return result;
  },
  view: function(vnode) {
    return m("main.layout", [
      m("h1", "Federation Heart Beat Information"),
      m("table", Discovery.dataIterator())
    ])
  },
  routes: {
    "/heartbeat" : {
      view: function() {
        return m(layout, m(Discovery))
      }
    }
  }
}

module.exports = Discovery
