import Ember from 'ember';

export default Ember.Route.extend({
  titleToken: function (page) {
    return page.get('name');
  },

  model: function (params) {
    return this.store.find('page', params.page_id);
  },

  actions: {
    error: function () {
      this.transitionTo('/not-found');
    },
  },
});
