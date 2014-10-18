export default Backbone.View.extend({
  template: _.template($('#template-page').html()),

  render: function () {
    this.$el.html(this.model.get("html"));
  },

  setModel: function (model) {
    this.model = model;
    if (!this.model.get("loaded")) {
      this.model.set('html', this.template(this.model.attributes));
      this.model.set("loaded", true);
      this.render();
    } else {
      this.render();
    }
  }
});
