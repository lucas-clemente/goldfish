import Page from 'page';

export default Backbone.Model.extend({
  namespace: '/v1',

  defaults: {
    pages: [],
    subFolders: [],
  },

  sync: function (method, collection, opts) {
    var _this = this;
    return $.ajax(this.namespace + this.id)
      .done(function (data) {
        var pages = data.map(function (id) {
          return new Page({id: id});
        });
        opts.success({pages: pages});
      })
      .fail(opts.fail);
  }
});
