import Ember from 'ember';

export default Ember.Route.extend({
  titleToken: function (folder) {
    return folder.get('path');
  },

  model: function (params) {
    return this.store.find('folder', params.folder_id);
  },
});
