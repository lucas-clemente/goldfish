export default Backbone.View.extend({
  tagName: 'a',
  className: 'list-group-item folder-item',
  template: _.template($('#template-folder-item').html()),

  render: function () {
    var href = this.model.id;
    this.$el.attr('href', href);
    this.$el.html(this.template(this.model));
    return this;
  },
});
