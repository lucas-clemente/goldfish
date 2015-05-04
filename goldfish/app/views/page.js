import Ember from 'ember';

export default Ember.View.extend({
  didInsertElement: function () {
    this.runTex();
  },

  runTex: Ember.observer('controller.model.compiled', function () {
    this.$('script').each(function (i, e) {
      var t = e.getAttribute('type');
      if (t.search('math/tex') === 0) {
        katex.render(e.textContent,
                     e.parentElement,
                     {
                       displayMode: t.search('mode=display') !== -1
                     });
      }
    });
  }),
});
