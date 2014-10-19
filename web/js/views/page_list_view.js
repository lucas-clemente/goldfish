import PageItemView from 'page_item_view';

export default Backbone.View.extend({
  tagName: "ul",
  className: 'list-group',

  initialize: function () {
    this.model.on("change", this.render, this);
  },

  render: function () {
    this.$el.empty();
    this.model.attributes.pages.forEach(function (page) {
      var itemView = new PageItemView({model: page});
      this.$el.append(itemView.render().el);
    }, this);
  },
});
