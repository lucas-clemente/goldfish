import Page from 'page';

var PageItemView = Backbone.View.extend({
  tagName: 'a',
  className: 'list-group-item',

  render: function () {
    this.$el.text(this.model.id);
    this.$el.attr('href', this.model.id);
    return this;
  },
});

var PageList = Backbone.Collection.extend({
  model: Page
});

var PageListView = Backbone.View.extend({
  tagName: "ul",
  className: 'list-group',

  initialize: function () {
    this.collection.on("reset", this.render, this);
  },

  render: function () {
    this.addAll();
  },

  addOne: function (page) {
    var itemView = new PageItemView({model: page});
    this.$el.append(itemView.render().el);
  },

  addAll: function () {
    this.collection.forEach(this.addOne, this);
  }
});

var PageView = Backbone.View.extend({
  render: function () {
    this.$el.html(this.model.get("html"));
  },

  setModel: function (model) {
    this.model = model;
    if (!this.model.get("loaded")) {
      this.model.set("html", "<h1>" + this.model.get("id") + "</h1>");
      this.model.set("loaded", true);
      this.render();
    } else {
      this.render();
    }
  }
});


var App = Backbone.Router.extend({
  routes: {
    "":              "folder",
    "*folder/":      "folder",
    "*folder/:page": "page",
    ":page":         "page",
  },

  folder: function (folder) {
    console.log("folder: ", folder);
  },

  page: function (folder, page) {
    if (!page) {
      page = folder;
      folder = '';
    }
    var model = window.appView.pageList.get(folder + '/' + page);
    window.appView.pageView.setModel(model);
  },
});


var AppView = Backbone.View.extend({
  initialize: function () {
    this.pageView = new PageView({model: new Page({id: "/"})});
    $('#page').append(this.pageView.el);
    this.pageList = new PageList();
    this.pageListView = new PageListView({collection: this.pageList});
    $('#list').append(this.pageListView.el);
    this.fetchCollections();
  },

  fetchCollections: function () {
    var items = [
      {id: "/foo", text: "item one"},
      {id: "/bar", text: "item two"},
      {id: "/baz", text: "item three"},
    ];
    this.pageList.reset(items);
  }
});

// From https://gist.github.com/tbranyen/1142129
$(document).delegate("a", "click", function(evt) {
  // Get the anchor href and protcol
  var href = $(this).attr("href");
  var protocol = this.protocol + "//";

  // Ensure the protocol is not part of URL, meaning its relative.
  // Stop the event bubbling to ensure the link will not cause a page refresh.
  if (href.slice(protocol.length) !== protocol) {
    evt.preventDefault();
    window.app.navigate(href, {trigger: true});
  }
});


$(function () {
  window.app = new App();
  window.appView = new AppView();
  Backbone.history.start({pushState: true});
});
