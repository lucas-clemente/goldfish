import Ember from 'ember';

export default Ember.Controller.extend({
  isEditing: false,
  oldSource: null,

  actions: {
    startEditing: function () {
      this.oldSource = this.get('model.markdownSource');
      this.toggleProperty('isEditing');
    },

    discardChanges: function () {
      this.set('model.markdownSource', this.oldSource);
      this.toggleProperty('isEditing');
    },

    saveChanges: function () {
      this.get('model').saveMarkdown();
      this.toggleProperty('isEditing');
    },
  },
});
