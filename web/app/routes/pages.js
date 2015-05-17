import Ember from 'ember';

export default Ember.Route.extend({
  titleToken: 'Pages',

  queryParams: {
    q: {
      refreshModel: true,
      replace: true,
    },
  },

  model: function (params) {
    return this.store.find('page', params);
  },
});
