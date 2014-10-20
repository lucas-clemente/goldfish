import fileFactory from 'file_factory';

var Folder = Backbone.Model.extend({
  namespace: '/v1',

  defaults: {
    files: [],
    subFolders: [],
  },

  name: function () {
    return this.id.match(/\/([^\/]+)\/$/)[1];
  },

  sync: function (method, collection, opts) {
    var _this = this;
    return $.ajax(this.namespace + this.id)
      .done(function (data) {
        var files = [], subFolders = [];
        data.forEach(function (path) {
          if (path[path.length-1] == '/') {
            subFolders.push(new Folder({id: path}));
          } else {
            var klass = fileFactory(path);
            files.push(new klass({id: path}));
          }
        });
        opts.success({
          files: files,
          subFolders: subFolders,
        });
      })
      .fail(opts.fail);
  },
});

export default Folder;
