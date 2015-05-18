import Ember from 'ember';

export default Ember.ArrayController.extend({
  queryParams: ['q'],
  q: null,

  sortProperties: ['modifiedAt'],
  sortAscending: false,
});
