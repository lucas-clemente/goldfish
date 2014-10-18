import Page from 'page';
import PageListView from 'page_list_view';
import PageView from 'page_view';

var PageList = Backbone.Collection.extend({
  model: Page
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
      {id: "/foo", text: "Foobar _emph_\n\n __strong__\n\n- list 1\n- list 2"},
      {id: "/bar", text: "item two"},
      {id: "/baz", text: "item three"},
    ];
    this.pageList.reset(items);
  }
});

// From https://gist.github.com/tbranyen/1142129
$(document).delegate("a", "click", function(e) {
  var href = $(this).attr("href");
  var protocol = this.protocol + "//";

  if (href.slice(protocol.length) !== protocol) {
    e.preventDefault();
    window.app.navigate(href, {trigger: true});
  }
});

$(function () {
  window.app = new App();
  window.appView = new AppView();
  Backbone.history.start({pushState: true});
});
