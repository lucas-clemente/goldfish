var Page = Backbone.Model.extend({
  namespace: '/v1',
  icon: 'fa-file-text-o',
  type: "page",

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
    var url = 'ws://' + window.location.host + this.namespace + this.id;
    console.log('connecting to', url);
    this.ws = new WebSocket(url);
    this.ws.onmessage = function (event) {
      opts.success({text: event.data});
    };
    this.ws.onerror = function (event) {
      opts.fail();
    };
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
    window.app.setTitle(this.attributes.title);

    // Escape \[ and \] to \\[ and \\]
    markdownRaw = markdownRaw.replace(/\\\[/g, '\\\\[');
    markdownRaw = markdownRaw.replace(/\\\]/g, '\\\\]');

    // Replace [[links]]

    // Images, e.g. [[foo.jpg]]
    markdownRaw = markdownRaw.replace(/\[\[([^|\]]+\.(?:jpg|png))\]\]/g, "![$1]($1)");
    // [[foo.extension]]
    markdownRaw = markdownRaw.replace(/\[\[([^|\]]+)\.([^|\]\.]+)\]\]/g, "[$1]($1.$2)");
    // [[foo]]
    markdownRaw = markdownRaw.replace(/\[\[([^|\]]+)\]\]/g, "[$1]($1.md)");
    // [[foo|bar.extension]]
    markdownRaw = markdownRaw.replace(/\[\[([^|\]]+)\|([^\]]+)\.([^\]\.]+)\]\]/g, "[$1]($2.$3)");
    markdownRaw = markdownRaw.replace(/\[\[([^|\]]+)\|([^\]]+)\]\]/g, "[$1]($2.md)");

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

  release: function () {
    this.ws.close();
  },
});

export default Page;
