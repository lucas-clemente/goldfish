export default Backbone.Model.extend({
  defaults: {
    html: "",
    loaded: false,
    text: ""
  },

  markdownText: function () {
    return marked(this.attributes.text);
  },

  name: function () {
    return this.id.slice(this.id.lastIndexOf('/') + 1);
  },
});
