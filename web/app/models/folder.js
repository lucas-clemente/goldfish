import Ember from 'ember';
import DS from 'ember-data';

export default DS.Model.extend({
  pages: DS.hasMany('page', {async: true}),
  subfolders: DS.hasMany('folder', {async: true, inverse: 'parentFolder'}),
  parentFolder: DS.belongsTo('folder', {async: true}),

  name: Ember.computed('id', function () {
    var name = this.id.slice(this.id.lastIndexOf('|') + 1);
    if (name === "") {
      name = "/";
    }
    return name;
  }),

  path: Ember.computed('id', function () {
    return this.id.replace(/\|/g, '/');
  })
});
