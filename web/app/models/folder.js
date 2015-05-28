import Ember from 'ember';
import DS from 'ember-data';

export default DS.Model.extend({
  pages: DS.hasMany('page', {async: true}),
  subfolders: DS.hasMany('folder', {async: true, inverse: 'parentFolder'}),
  parentFolder: DS.belongsTo('folder', {async: true}),

  sortedPages: Ember.computed.sort('pages', function (a, b) {
    if (a.id < b.id) {
      return -1;
    } else if (a.id > b.id) {
      return 1;
    }
    return 0;
  }),

  name: Ember.computed('id', function () {
    var name = this.id.slice(this.id.lastIndexOf('|') + 1);
    if (name === "") {
      name = "/";
    }
    return name;
  }),

  path: Ember.computed('id', function () {
    return this.id.replace(/\|/g, '/');
  }),

  isRootAndEmpty: Ember.computed('subfolders.[]', 'pages.[]', 'id', function () {
    return this.get('subfolders.length') === 0 && this.get('pages.length') === 0 && this.id === '|';
  }),
});
