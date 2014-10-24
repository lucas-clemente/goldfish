import Page from 'page';
import Folder from 'folder';
import FolderListView from 'folder_list_view';
import DetailView from 'detail_view';
import fileFactory from 'file_factory';


var App = Backbone.Router.extend({
  routes: {
    "":              "folder",
    "*folder/":      "folder",
    "*folder/:file": "detail",
    ":file":         "detail",
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
    this.setTitle(folder);
    window.appView.setFolder(folder);
  },

  detail: function (folder, path) {
    if (!path) {
      // Root path, swap path and folder
      folder = [path, path = folder][0];
      path = "/" + path;
    } else {
      path = "/" + folder + '/' + path;
    }
    console.log("setting detail path", path);
    this.folder(folder);
    this.setTitle(path.slice(path.lastIndexOf('/') + 1));
    window.appView.setDetail(path);
  },

  setTitle: function (title) {
    document.title = title + " — Goldfish";
  },
});


var AppView = Backbone.View.extend({
  folder: new Folder(),

  initialize: function () {
    this.folderListView = new FolderListView({model: this.folder});
    $('#list').append(this.folderListView.el);
    this.detailView = new DetailView();
    $('#page').append(this.detailView.el);
  },

  setFolder: function (path) {
    this.folder.id = path;
    this.folder.fetch();
  },

  setDetail: function (path) {
    var klass = fileFactory(path);
    var model = new klass({id: path});
    var _this = this;
    model.fetch({
      success: function () {
        _this.detailView.setModel(model);
      },
      fail: function () {
        _this.detailView.set404(path);
      }
    });
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

window.MathJax = {
  tex2jax: {
    inlineMath: [['$$', '$$']],
    displayMath: [['\\[', '\\]']],
  },
  "HTML-CSS": {
    scale: 90,
    preferredFont: "STIX-Web",
    webFont: "STIX-Web",
    imageFont: null,
  }
};

$(function () {
  window.app = new App();
  window.appView = new AppView();
  Backbone.history.start({pushState: true});
});
