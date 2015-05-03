import Ember from 'ember';
import DS from 'ember-data';

export default DS.Model.extend({
  path: DS.attr('string'),
  text: DS.attr('string'),

  title: Ember.computed('path', 'text', function () {
    var path = this.get('path');
    return path.slice(path.lastIndexOf('/') + 1);
  }),
});
