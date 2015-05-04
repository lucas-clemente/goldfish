import Ember from 'ember';
import DS from 'ember-data';

export default DS.Model.extend({
  path: DS.attr('string'),
  text: DS.attr('string'),

  markdownRenderer: null,

  // Either the top level heading or the filename
  title: Ember.computed('path', 'text', function () {
    var m = /^#(.*)\n/.exec(this.get('text'));
    if (m) {
      return m[1].trim();
    }
    var path = this.get('path');
    return path.slice(path.lastIndexOf('/') + 1);
  }),

  init: function () {
    this._super.apply(this, arguments);
    this.markdownRenderer = new marked.Renderer();

    this.markdownRenderer.image = (href, title) => {
      if (href[0] === '/') {
        href = '/v1' + href;
      } else {
        href = '/v1' + this.currentFolder() + href;
      }
      return '<div class="image"><img src="' + href + '" title="' + (title || '') + '" class="img-thumbnail" /></div>';
    };

    // Make all links absolute
    this.markdownRenderer.link = (href, title, text) => {
      if (href.search(/http/) !== 0 && href[0] !== '/') {
        href = this.currentFolder() + href;
      }
      return '<a href="' + href + '">' + text + '</a>';
    };
  },

  compiled: Ember.computed('text', function () {
    var source = this.get('text') || "";

    // Remove top level heading
    source = source.replace(/^#(.*)/, "");

    // Replace [[links]]

    // Images, e.g. [[foo.jpg]]
    source = source.replace(/\[\[([^|\]]+\.(?:jpg|png))\]\]/g, "![$1]($1)");
    // [[foo.extension]]
    source = source.replace(/\[\[([^|\]]+)\.([^|\]\.]+)\]\]/g, "[$1]($1.$2)");
    // [[foo]]
    source = source.replace(/\[\[([^|\]]+)\]\]/g, "[$1]($1.md)");
    // [[foo|bar.extension]]
    source = source.replace(/\[\[([^|\]]+)\|([^\]]+)\.([^\]\.]+)\]\]/g, "[$1]($2.$3)");
    source = source.replace(/\[\[([^|\]]+)\|([^\]]+)\]\]/g, "[$1]($2.md)");


    var compiled = marked(source, {renderer: this.markdownRenderer});

    return compiled;
  }),

  currentFolder: function () {
    return "/" + this.id.slice(0, this.id.lastIndexOf('/')+1);
  },
});
