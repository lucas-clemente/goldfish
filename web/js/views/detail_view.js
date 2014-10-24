var DetailView = Backbone.View.extend({
  templates: {
    page: _.template($('#template-page').html()),
    image: _.template($('#template-image').html()),
    file: _.template($('#template-file').html()),
    notFound: _.template($('#template-notfound').html()),
  },

  initialize: function () {
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
    if (this.model && this.model !== model && this.model.release) {
      this.model.release();
    }
    this.model = model;
    this.render();
  },

  set404: function (path) {
    this.model = {type: "notFound", id: path};
    this.render();
  },
});

export default DetailView;
