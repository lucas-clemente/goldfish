import PageItemView from 'page_item_view';
import FolderItemView from 'folder_item_view';

export default Backbone.View.extend({
  tagName: "ul",
  className: 'list-group',

  initialize: function () {
    this.model.on("change", this.render, this);
  },

  render: function () {
    this.$el.empty();
    this.model.attributes.subFolders.forEach(function (f) {
      this.$el.append((new FolderItemView({model: f})).render().el);
    }, this);
    this.model.attributes.pages.forEach(function (p) {
      this.$el.append((new PageItemView({model: p})).render().el);
    }, this);
  },
});
