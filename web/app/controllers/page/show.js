import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
    deletePage: function () {
      if (!window.confirm('Do you really want to delete this page?')) {
        return;
      }
      var folder = this.get('model.folder');
      this.get('model').destroyRecord();
      this.transitionToRoute('folder', folder);
    },
  },
});
