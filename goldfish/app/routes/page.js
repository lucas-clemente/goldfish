import Ember from 'ember';

export default Ember.Route.extend({
  model: function (params) {
    var id = params.path;
    var folder = id.slice(0, id.lastIndexOf('/')+1);
    return new Ember.RSVP.Promise((resolve, reject) => {
      Ember.$.get('/v1/' + id)
      .done((text) => {
        var page = this.store.push('page', {
          id: id,
          folder: folder,
          text: text,
        });
        resolve(page);
      })
      .fail(reject);
    });
  }
});
