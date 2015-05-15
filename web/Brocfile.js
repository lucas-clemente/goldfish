/* global require, module */

var EmberApp = require('ember-cli/lib/broccoli/ember-app');
var glob = require('glob');

var app = new EmberApp();

// Use `app.import` to add additional libraries to the generated
// output files.
//
// If you need to use different assets in different
// environments, specify an object as the first parameter. That
// object's keys should be the environment name and the values
// should be the asset to use in that environment.
//
// If the library that you are including contains AMD or ES6
// modules that you would like to import into your application
// please specify an object with the list of modules as keys
// along with the exports of each module as its value.

app.import('bower_components/bootstrap-sass-official/assets/javascripts/bootstrap.js');

var fontDir = {
  destDir: 'assets/fonts'
};

app.import('bower_components/fontawesome/fonts/fontawesome-webfont.woff', fontDir);
app.import('bower_components/fontawesome/fonts/fontawesome-webfont.woff2', fontDir);

app.import('bower_components/highlightjs/highlight.pack.js');
app.import('bower_components/highlightjs/styles/tomorrow.css');

app.import('bower_components/marked/lib/marked.js');

app.import('bower_components/katex-build/katex.min.js');
app.import('bower_components/katex-build/katex.min.css');

var fontFiles = glob.sync('bower_components/katex-build/fonts/*.woff?(2)');
for (var i = 0; i < fontFiles.length; i++) {
  app.import(fontFiles[i], fontDir);
}

app.import('bower_components/autosize/dest/autosize.js');

module.exports = app.toTree();
