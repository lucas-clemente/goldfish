import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
    discardChanges: function () {
      this.get('model').rollback();
      this.transitionToRoute('page.show', this.get('model'));
    },

    saveChanges: function () {
      this.get('model').save();
      this.transitionToRoute('page.show', this.get('model'));
    },
  },
});
