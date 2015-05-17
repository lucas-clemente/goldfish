import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
  location: config.locationType,
});

export default Router.map(function() {
  this.resource('folder', {path: '/folders/:folder_id'}, function () {
    this.resource('page', {path: '/pages/:page_id'});
  });
  this.route('not-found', {path: '/*path'});
  this.route('pages', {path: '/pages'});
});
