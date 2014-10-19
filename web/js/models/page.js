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
    this.attributes.markdown = marked(markdownRaw);
  },
});
