// The following is required to ensure that webpack v4.x.x can use
// YARN's advanced PlugNPlay package resolution
//
// see: https://github.com/arcanis/pnp-webpack-plugin
//

const PnpWebpackPlugin = require(`pnp-webpack-plugin`);

module.exports = {
  resolve: {
    plugins: [
      PnpWebpackPlugin,
    ],
  },
  resolveLoader: {
    plugins: [
      PnpWebpackPlugin.moduleLoader(module),
    ],
  },
};
