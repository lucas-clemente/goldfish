import Ember from 'ember';
import DS from 'ember-data';

export default DS.Model.extend({
  path: DS.attr('string'),
  text: DS.attr('string'),

  markdownRenderer: null,

  title: Ember.computed('path', 'text', function () {
    var path = this.get('path');
    return path.slice(path.lastIndexOf('/') + 1);
  }),

  init: function () {
    this._super.apply(this, arguments);
    this.markdownRenderer = new marked.Renderer();

    this.markdownRenderer.image = (href, title, text) => {
      if (href[0] === '/') {
        href = '/v1' + href;
      } else {
        href = '/v1/' + this.id.slice(0, this.id.lastIndexOf('/')+1) + href;
      }
      return '<div class="image"><img src="' + href + '" title="' + (title || '') + '" class="img-thumbnail" /></div>';
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
});
