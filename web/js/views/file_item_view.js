export default Backbone.View.extend({
  tagName: 'a',
  className: 'list-group-item file-item',
  template: _.template($('#template-file-item').html()),

  render: function () {
    var href = this.model.id;
    this.$el.attr('href', href);
    this.$el.html(this.template(this.model));
    if (decodeURIComponent(window.location.pathname) === href) {
      this.$el.addClass('active');
    }
    return this;
  },
});
