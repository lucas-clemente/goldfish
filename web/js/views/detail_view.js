var DetailView = Backbone.View.extend({
  templates: {
    page: _.template($('#template-page').html()),
    image: _.template($('#template-image').html()),
    file: _.template($('#template-file').html()),
  },

  initialize: function () {
    // this.model.on('change', this.render, this);
  },

  render: function () {
    var template = this.templates[this.model.type];
    if (!template) {
      console.error('cannot find template for', this.model.type);
    }
    this.$el.html(template(this.model));

    // Typeset latex
    if (this.model.type === "page") {
      MathJax.Hub.Queue(['Typeset', MathJax.Hub]);
    }
  },

  setModel: function (model) {
    if (this.model) {
      this.model.off('change', this.render, this);
    }
    this.model = model;
    if (!this.model.loading) {
      this.render();
    }
    this.model.on('change', this.render, this);
  },
});

export default DetailView;
