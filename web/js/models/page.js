export default Backbone.Model.extend({
  namespace: '/v1',

  defaults: {
    text: "",
    markdown: "",
    title: "",
  },

  initialize: function () {
    this.on("change", this.updateMarkdown, this);
  },

  name: function () {
    return this.id.slice(this.id.lastIndexOf('/') + 1);
  },

  sync: function (method, collection, opts) {
    var _this = this;
    return $.ajax(this.namespace + this.id)
      .done(function (data) {
        opts.success({text: data});
      })
      .fail(opts.fail);
  },

  updateMarkdown: function () {
    // Take markdown title if matching
    var m = /^#(.*)\n([^]*)$/.exec(this.attributes.text);
    var markdownRaw = "";
    if (m) {
      this.attributes.title = m[1].trim();
      markdownRaw = m[2];
    } else {
      this.attributes.title = this.name();
      markdownRaw = this.attributes.text;
    }

    // Escape \[ and \] to \\[ and \\]
    markdownRaw = markdownRaw.replace('\\[', '\\\\[');
    markdownRaw = markdownRaw.replace('\\]', '\\\\]');

    // Replace [[links]]

    // Images, e.g. [[foo.jpg]]
    markdownRaw = markdownRaw.replace(/\[\[([^|\]]+\.(?:jpg|png))\]\]/g, "![$1]($1)");
    // [[foo]]
    markdownRaw = markdownRaw.replace(/\[\[([^|\]]+)\]\]/g, "[$1]($1)");
    // [[foo|bar]]
    markdownRaw = markdownRaw.replace(/\[\[([^|\]]+)\|([^\]]+)\]\]/g, "[$1]($2)");

    // Render markdown

    var renderer = new marked.Renderer();
    var _this = this;
    renderer.image = function (href, title, text) {
      if (href[0] == '/') {
        href = '/v1' + href;
      } else {
        href = '/v1' + _this.id.slice(0, _this.id.lastIndexOf('/')+1) + href;
      }
      return '<div class="image"><img src="' + href + '" title="' + (title || '') + '" class="img-thumbnail" /></div>';
    };
    this.attributes.markdown = marked(markdownRaw, {renderer: renderer});
  },
});
