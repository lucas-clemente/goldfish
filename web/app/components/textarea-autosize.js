import Ember from 'ember';

export default Ember.TextArea.extend({
  didInsertElement: function () {
    this._super();
    Ember.run.scheduleOnce('afterRender', this, function () {
      autosize(this.$());
    });
  },

  willDestroyElement: function () {
    this._super();
    autosize.destroy(this.$());
  },
});
