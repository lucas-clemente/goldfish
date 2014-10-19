import Page from 'page';
import Folder from 'folder';
import PageListView from 'page_list_view';
import PageView from 'page_view';


var App = Backbone.Router.extend({
  initialize: function () {
    this.bind('all', this.updateHighlights);
  },

  routes: {
    "":              "folder",
    "*folder/":      "folder",
    "*folder/:page": "page",
    ":page":         "page",
  },

  updateHighlights: function (route) {
    $('a[href="' + window.location.pathname + '"]').addClass("active");
    $('a[href!="' + window.location.pathname + '"]').removeClass("active");
  },

  folder: function (folder) {
    window.appView.setFolder(folder || '/');
  },

  page: function (folder, page) {
    if (!page) {
      // Root page
      page = folder;
      folder = '/';
    } else {
      page = folder + '/' + page;
    }
    window.appView.setFolder(folder);
    var model = window.appView.pageList.get(page);
    window.appView.pageView.setModel(model);
  },
});


var AppView = Backbone.View.extend({
  folder: new Folder(),

  initialize: function () {
    this.pageView = new PageView({model: new Page({id: "/"})});
    $('#page').append(this.pageView.el);
    this.pageListView = new PageListView({model: this.folder});
    $('#list').append(this.pageListView.el);
  },

  setFolder: function (path) {
    this.folder.id = path;
    var _this = this;
    this.folder.fetch();
  },
});


// From https://gist.github.com/tbranyen/1142129
$(document).delegate("a", "click", function(e) {
  var href = this.getAttribute("href");
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
