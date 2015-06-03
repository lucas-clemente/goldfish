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

    uploadAndLinkFile: function (fileList, textArea) {
      var textToInsert = '';

      var folder = this.get('page.folder');
      for (var i = 0; i < fileList.length; i++) {
        var file = fileList.item(i);
        var id = folder.get('id') + '|' + file.name;
        id = id.replace('||', '|');
        var page = this.store.createRecord('page', {
          id: id,
          folder: folder,
        });
        /* jshint -W083 */
        page.save().then(function () {
          page.sendData(file);
        });

        textToInsert += '[[' + file.name + ']]';
      }

      var pos = textArea.selectionStart || 0;
      var text = this.get('page.markdownSource');
      this.set('page.markdownSource', text.substring(0, pos) + textToInsert + text.substring(pos));
      textArea.selectionStart = pos + textToInsert.length;
      textArea.selectionEnd = pos + textToInsert.length;
    },
  },
});
