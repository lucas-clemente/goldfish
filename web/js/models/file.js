var File = Backbone.Model.extend({
  namespace: '/v1',
  icon: 'fa-file-o',
  type: "file",

  name: function () {
    return this.id.slice(this.id.lastIndexOf('/') + 1);
  },

  sync: function (method, collection, opts) {
    var d = $.Deferred();
    d.resolve();
    return d;
  },
});

export default File;
