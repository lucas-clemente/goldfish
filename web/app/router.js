import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.resource('folder', {path: '/folders/:folder_id'}, function () {
    this.resource('page', {path: '/pages/:page_id'}, function () {
      this.route('show', {path: ''});
      this.route('edit');
    });
  });
  this.route('not-found', {path: '/*path'});
  this.route('pages', {path: '/pages'});
});

export default Router;
