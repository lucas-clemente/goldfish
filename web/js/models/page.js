export default Backbone.Model.extend({
  defaults: {
    html: "",
    loaded: false,
    text: ""
  },

  markdownText: function () {
    return marked(this.attributes.text);
  }
});
