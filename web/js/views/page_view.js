export default Backbone.View.extend({
  template: _.template($('#template-page').html()),

  initialize: function () {
    this.model.on("change", this.render, this);
  },

  render: function () {
    this.$el.html(this.template(this.model.attributes));
    MathJax.Hub.Queue(["Typeset", MathJax.Hub]);
  },
});
