import Ember from 'ember';

export default Ember.Controller.extend({
  addingNewFile: false,
  newFileName: '',

  actions: {
    startStopAddingNewFile: function () {
      this.toggleProperty('addingNewFile');
    },

    addNewFile: function () {
      var filename = this.get('newFileName') + '.md';
      var folder = this.get('model');

      var file = this.store.createRecord('page', {
        id: folder.id + '|' + filename,
        folder: folder,
      });

      this.transitionTo('page.edit', file);

      this.set('newFileName', '');
      this.set('addingNewFile', false);
    },
  },
});
