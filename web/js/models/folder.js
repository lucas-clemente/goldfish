import Page from 'page';

var Folder = Backbone.Model.extend({
  namespace: '/v1',

  defaults: {
    pages: [],
    subFolders: [],
  },

  name: function () {
    return this.id.match(/\/([^\/]+)\/$/)[1];
  },

  sync: function (method, collection, opts) {
    var _this = this;
    return $.ajax(this.namespace + this.id)
      .done(function (data) {
        var pages = data
          .filter(function (id) {
            return id[id.length-1] !== "/";
          })
          .map(function (id) {
            return new Page({id: id});
          });
        var subFolders = data
          .filter(function (id) {
            return id[id.length-1] === "/";
          })
          .map(function (id) {
            return new Folder({id: id});
          });
        opts.success({
          pages: pages,
          subFolders: subFolders,
        });
      })
      .fail(opts.fail);
  },
});

export default Folder;
