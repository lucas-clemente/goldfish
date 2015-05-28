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

      var id = folder.id + '|' + filename;
      id = id.replace('||', '|');

      var file = this.store.createRecord('page', {
        id: id,
        folder: folder,
      });

      this.transitionToRoute('page.edit', file);

      this.set('newFileName', '');
      this.set('addingNewFile', false);
    },

    uploadFiles: function (fileList) {
      var folder = this.get('model');
      for (var i = 0; i < fileList.length; i++) {
        var file = fileList.item(i);
        var id = folder.id + '|' + file.name;
        id = id.replace('||', '|');
        var page = this.store.createRecord('page', {
          id: id,
          folder: folder,
        });
        page.sendData(file);
      }
    },
  },
});
