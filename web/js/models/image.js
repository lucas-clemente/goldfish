var Image = Backbone.Model.extend({
  namespace: '/v1',
  icon: 'fa-file-image-o',
  type: "image",

  name: function () {
    return this.id.slice(this.id.lastIndexOf('/') + 1);
  },

  sync: function (method, collection, opts) {
    var d = $.Deferred();
    d.resolve();
    return d;
  },
});

export default Image;
