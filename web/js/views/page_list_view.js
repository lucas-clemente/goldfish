import PageItemView from 'page_item_view';

export default Backbone.View.extend({
  tagName: "ul",
  className: 'list-group',

  initialize: function () {
    this.collection.on("reset", this.render, this);
  },

  render: function () {
    this.addAll();
  },

  addOne: function (page) {
    var itemView = new PageItemView({model: page});
    this.$el.append(itemView.render().el);
  },

  addAll: function () {
    this.collection.forEach(this.addOne, this);
  }
});
