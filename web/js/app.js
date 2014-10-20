import Page from 'page';
import Folder from 'folder';
import FolderListView from 'folder_list_view';
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
    if (folder) {
      // Normal folder
      folder = "/" + folder + "/";
    } else {
      // Root
      folder = "/";
    }
    console.log("setting folder", folder);
    window.appView.setFolder(folder);
  },

  page: function (folder, page) {
    if (!page) {
      // Root page, swap page and folder
      folder = [page, page = folder][0];
      page = "/" + page;
    } else {
      page = "/" + folder + '/' + page;
    }
    console.log("setting page", page);
    this.folder(folder);
    window.appView.setPage(page);
  },
});


var AppView = Backbone.View.extend({
  folder: new Folder(),
  page: new Page(),

  initialize: function () {
    this.folderListView = new FolderListView({model: this.folder});
    $('#list').append(this.folderListView.el);
    this.pageView = new PageView({model: this.page});
    $('#page').append(this.pageView.el);
  },

  setFolder: function (path) {
    this.folder.id = path;
    this.folder.fetch();
  },

  setPage: function (path) {
    this.page.id = path;
    this.page.fetch();
  },
});


$(document).on('click', 'a', function(e) {
  var href = this.getAttribute("href");
  if (href[0] == "/") {
    e.preventDefault();
    window.app.navigate(href, {trigger: true});
  }
});

marked.setOptions({
  highlight: function (code, lang) {
    if (lang) {
      return hljs.highlight(lang, code).value;
    } else {
      return hljs.highlightAuto(code).value;
    }
  }
});

$(function () {
  window.app = new App();
  window.appView = new AppView();
  Backbone.history.start({pushState: true});
});
