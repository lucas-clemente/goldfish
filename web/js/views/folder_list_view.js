import FileItemView from 'file_item_view';
import FolderItemView from 'folder_item_view';

export default Backbone.View.extend({
  tagName: "ul",
  className: 'list-group',
  currentTemplate: _.template($("#template-current-folder").html()),
  folderUpTemplate: _.template($("#template-folder-up").html()),

  initialize: function () {
    this.model.on("change", this.render, this);
  },

  render: function () {
    this.$el.empty();
    this.$el.append($(this.currentTemplate({url: this.model.id})));
    if (this.model.id != "/") {
      var upURL = this.model.id.slice(0, this.model.id.slice(0, -1).lastIndexOf("/")+1);
      this.$el.append($(this.folderUpTemplate({url: upURL})));
    }
    this.model.attributes.subFolders.forEach(function (f) {
      this.$el.append((new FolderItemView({model: f})).render().el);
    }, this);
    this.model.attributes.files.forEach(function (p) {
      this.$el.append((new FileItemView({model: p})).render().el);
    }, this);
  },
});
