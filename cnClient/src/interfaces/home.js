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

var Home = {
  data: {
    "Desc" : "",
    "SubRoutes": []
  },
  oninit: function() {
    m.request({
      method: "GET",
      url: "/",
    }).then(function(result){
      console.dir(result, {depth: null, colors: true});
      if ( result == null ) {
        result = {
          "Desc" : "",
          "SubRoutes" : []
        };
      }
      Home.data = result;
    })
  },
  view: function(vnode) {
    return m("main.layout", [
      m("p", Home.data.Desc),
      m("ul", 
        Home.data.SubRoutes.map(function(aSubRoute) {
          return aSubRoute.Visible && m("li",
            m("strong", 
              m(m.route.Link,
                { href: aSubRoute.Path },
                aSubRoute.Path
              )
            ),
            m("p", aSubRoute.Desc)
          )
        })
      )
    ])
  },
  routes: {
    "/home" : {
      view: function() {
        return m(layout, m(Home))
      }
    }
  }
}

module.exports = Home
