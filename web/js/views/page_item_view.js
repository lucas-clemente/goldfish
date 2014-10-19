export default Backbone.View.extend({
  tagName: 'a',
  className: 'list-group-item file-item',
  template: _.template($('#template-page-item').html()),

  render: function () {
    this.$el.attr('href', this.model.id);
    this.$el.html(this.template(this.model));
    return this;
  },
});
