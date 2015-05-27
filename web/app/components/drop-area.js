import Ember from 'ember';

export default Ember.Component.extend({
  classNameBindings: ['dragging'],
  dragging: false,

  didInsertElement: function () {
    var el = this.$()[0];
    var counter = 0; // Counter so that entering a child doesn't stop dragging.
    el.addEventListener('dragenter', (e) => {
      if (e.dataTransfer.files.length === 0) {
        return;
      }
      counter++;
      e.preventDefault();
      e.dataTransfer.effectAllowed = 'copy';
      e.dataTransfer.dropEffect = 'copy';
      this.set('dragging', true);
    });
    el.addEventListener('dragover', (e) => {
      e.preventDefault();
    });
    el.addEventListener('dragleave', (e) => {
      if (e.dataTransfer.files.length === 0) {
        return;
      }
      counter--;
      if (counter !== 0) {
        return;
      }
      e.preventDefault();
      this.set('dragging', false);
    });
    el.addEventListener('drop', (e) => {
      if (e.dataTransfer.files.length === 0) {
        return;
      }
      e.preventDefault();
      this.set('dragging', false);
      var files = e.dataTransfer.files;
      this.sendAction('action', files);
    });
  },
});
