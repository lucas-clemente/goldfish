import DS from 'ember-data';

export default DS.RESTAdapter.extend({
  namespace: 'v2',

  buildURL: function (type, id) {
    return this.namespace + '/' + type + 's/' + id;
  },
});
