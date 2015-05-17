import Ember from 'ember';

export default Ember.Controller.extend({
  searchText: '',

  actions: {
    'search': function () {
      this.transitionToRoute('pages', {queryParams: {q: this.get('searchText')}});
    },
  },
});
