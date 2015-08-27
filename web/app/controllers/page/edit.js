import Ember from 'ember';
import _function from 'lodash/function';

export default Ember.Controller.extend({
  saveDisallowed: Ember.computed.not('model.hasDirtyAttributes'),

  autosaveModel: Ember.observer('model.markdownSource', function() {
    this.debouncedAutosave();
  }),

  debouncedAutosave: _function.debounce(function() {
    if (this.get('model.hasDirtyAttributes')) {
      this.get('model').save();
    }
  }, 1000),

  actions: {
    finishEditing: function () {
      this.get('model').save();
      this.transitionToRoute('page.show', this.get('model'));
    },

    uploadAndLinkFile: function (fileList, textArea) {
      var textToInsert = '';

      var folder = this.get('model.folder');
      for (var i = 0; i < fileList.length; i++) {
        var file = fileList.item(i);
        var id = folder.get('id') + '|' + file.name;
        id = id.replace('||', '|');
        var model = this.store.createRecord('model', {
          id: id,
          folder: folder,
        });
        /* jshint -W083 */
        model.save().then(function () {
          model.sendData(file);
        });

        textToInsert += '[[' + file.name + ']]';
      }

      var pos = textArea.selectionStart || 0;
      var text = this.get('model.markdownSource');
      this.set('model.markdownSource', text.substring(0, pos) + textToInsert + text.substring(pos));
      textArea.selectionStart = pos + textToInsert.length;
      textArea.selectionEnd = pos + textToInsert.length;
    },
  },
});
