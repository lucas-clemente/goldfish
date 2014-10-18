export default Backbone.View.extend({
  tagName: 'a',
  className: 'list-group-item',

  render: function () {
    this.$el.text(this.model.id);
    this.$el.attr('href', this.model.id);
    return this;
  },
});
