import Ember from 'ember';

export default Ember.View.extend({
  didInsertElement: function () {
    this.updateDomStuff();
  },

  updateDomStuff: Ember.observer('controller.model.compiled', function () {
    Ember.run.scheduleOnce('afterRender', this, function() {
      // This is a workaround for #37.
      // For some reason this.$ is sometimes undefined.
      var el = this.get('element');
      if (!el) {
        return;
      }

      // Render tex
      Ember.$(el).find('script').each(function (i, e) {
        var t = e.getAttribute('type');
        if (t.search('math/tex') === 0) {
          katex.render(e.textContent,
                       e.parentElement,
                       {
                         displayMode: t.search('mode=display') !== -1
                       });
        }
      });

      // Fix internal links
      Ember.$(el).find('a').each((i, el) => {
        var href = el.getAttribute('href');
        if (href[0] === '/') {
          Ember.$(el).click((ev) => {
            var router = this.get('controller.target.router');
            router.transitionTo(href);
            ev.preventDefault();
          });
        }
      });
    });
  }),
});
