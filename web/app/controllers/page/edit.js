import Ember from 'ember';
import { computedAutosave } from 'ember-autosave';

export default Ember.Controller.extend({
  page: computedAutosave('model'),

  saveDisallowed: Ember.computed.not('page.isDirty'),

  actions: {
    finishEditing: function () {
      this.get('model').save();
      this.transitionToRoute('page.show', this.get('model'));
    },
  },
});
