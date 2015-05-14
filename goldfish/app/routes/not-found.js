import Ember from 'ember';

export default Ember.Route.extend({
  titleToken: '404',
  
  redirect: function () {
    var url = this.router.location.formatURL('/not-found');
    if (window.location.pathname !== url) {
      this.transitionTo('/not-found');
    }
  },
});
