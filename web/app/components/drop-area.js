import Ember from 'ember';

export default Ember.Component.extend({
  classNameBindings: ['dragging'],
  dragging: false,

  didInsertElement: function () {
    var el = this.$()[0];
    var counter = 0; // Counter so that entering a child doesn't stop dragging.
    el.addEventListener('dragenter', (e) => {
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
      counter--;
      if (counter !== 0) {
        return;
      }
      e.preventDefault();
      this.set('dragging', false);
    });
    el.addEventListener('drop', (e) => {
      e.preventDefault();
      this.set('dragging', false);
      if (e.dataTransfer.files.length === 0) {
        return;
      }
      var files = e.dataTransfer.files;
      this.sendAction('action', files);
    });
  },
});
