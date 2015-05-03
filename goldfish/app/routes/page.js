import Ember from 'ember';

export default Ember.Route.extend({
  model: function (params) {
    return new Ember.RSVP.Promise((resolve, reject) => {
      Ember.$.get('/v1/' + params.path)
      .done((text) => {
        var page = this.store.push('page', {
          id: params.path,
          path: params.path,
          text: text,
        });
        resolve(page);
      })
      .fail(reject);
    });
  }
});
