import Ember from 'ember';
import DS from 'ember-data';

export default DS.Model.extend({
  markdownSource: DS.attr('string'),
  modifiedAt: DS.attr('date'),
  folder: DS.belongsTo('folder', {async: true}),

  init: function () {
    this._super.apply(this, arguments);
    this.initMarkdownRenderer();
  },

  name: Ember.computed('id', function () {
    return this.id.slice(this.id.lastIndexOf('|') + 1);
  }),

  extension: Ember.computed('id', function () {
    return this.id.slice(this.id.lastIndexOf('.')+1);
  }),

  rawPath: Ember.computed('id', function () {
    return '/v2/raw' + this.id.replace(/\|/g, '/');
  }),

  icon: Ember.computed('extension', function () {
    switch (this.get('extension')) {
      case "md":
        return "file-text-o";
      case "png":
      case "jpg":
      case "svg":
        return "file-image-o";
      case "pdf":
        return "file-pdf-o";
      case "xls":
      case "xlsx":
        return "file-excel-o";
      case "ppt":
      case "pptx":
        return "file-powerpoint-o";
      case "doc":
      case "docx":
        return "file-word-o";
      default:
        return "file-o";
    }
  }),

  folderPath: Ember.computed('id', function () {
    var folderID = this.id.slice(0, this.id.lastIndexOf('|'));
    if (folderID === '') {
      folderID = '|';
    }
    return folderID.replace(/\|/g, '/');
  }),


  // Either the top level heading or the filename
  title: Ember.computed('id', 'markdownSource', function () {
    var m = /^#(.*)\n/.exec(this.get('markdownSource'));
    if (m) {
      return m[1].trim();
    }
    return this.get('name');
  }),

  open: function () {
    Ember.$.post('/v2/pages/' + this.id + '/open');
  },

  sendData: function (data) {
    var path = this.get('rawPath');

    var req = new XMLHttpRequest();
    req.onerror = function () {
      console.log('error uploading file');
      console.log(arguments);
    };
    req.open("POST", path, true);
    req.send(data);
    this.set('modifiedAt', new Date());
  },

  // --------------------------------------------------------------------------
  // -- Misc file formats -----------------------------------------------------
  // --------------------------------------------------------------------------

  isImage: Ember.computed('extension', function () {
    var ext = this.get('extension');
    return ext === 'jpg' || ext === 'png' || ext === 'svg';
  }),

  isPDF: Ember.computed('extension', function () {
    return this.get('extension') === 'pdf';
  }),

  // --------------------------------------------------------------------------
  // -- Markdown specific -----------------------------------------------------
  // --------------------------------------------------------------------------

  isMarkdown: Ember.computed('extension', function () {
    var ext = this.get('extension');
    return ext === 'md' || ext === 'markdown';
  }),

  initMarkdownRenderer: function () {
    this.markdownRenderer = new marked.Renderer();

    this.markdownRenderer.image = (href, title) => {
      if (href.search(/http/) !== 0) {
        if (href[0] === '/') {
          href = '/v2/raw' + href;
        } else {
          href = '/v2/raw' + this.get('folderPath') + "/" + href;
        }
      }
      return '<div class="image"><img src="' + href + '" title="' + (title || '') + '" class="img-thumbnail" /></div>';
    };

    // Make all links absolute
    this.markdownRenderer.link = (href, title, text) => {
      if (href.search(/http/) !== 0) {
        if (href[0] !== '/') {
          // Make link absolute
          href = this.get('folderPath') + '/' + href;
          href = href.replace('//', '/');
        }
        var refFolder = href.slice(0, href.lastIndexOf('/'));
        href = `/folders/${ refFolder.replace(/\//g, '|') }/pages/${ href.replace(/\//g, '|') }`;
      }
      return '<a href="' + href + '">' + text + '</a>';
    };

    marked.setOptions({
      highlight: function (code, lang) {
        if (lang) {
          return hljs.highlight(lang, code).value;
        } else {
          return hljs.highlightAuto(code).value;
        }
      }
    });
  },

  saveMarkdown: function () {
    var path = this.get('rawPath');
    Ember.$.post(path, this.get('markdownSource'))
    .fail(function () {
      console.error('error saving to ', path);
    });
    this.set('modifiedAt', new Date());
  },

  compiledMarkdown: Ember.computed('markdownSource', function () {
    var source = this.get('markdownSource');
    if (!source) {
      return "";
    }

    // Remove top level heading
    source = source.replace(/^#(.*)/, "");

    var multilineEqs = {};
    var inlineEqs = {};

    // Replace \[ and \] by placeholders
    source = source.replace(/\\\[([^]*?)\\\]/gm, function(m, eq) {
      var s = Math.random().toString(36).slice(2);
      multilineEqs[s] = eq;
      return s;
    });

    // Replace $ ... $ by placeholders
    source = source.replace(/\$([^]*?)\$/gm, function(m, eq) {
      var s = Math.random().toString(36).slice(2);
      inlineEqs[s] = eq;
      return s;
    });

    // Replace [[links]]

    // Images, e.g. [[foo.jpg]]
    source = source.replace(/\[\[([^|\]]+\.(?:jpg|png|svg))\]\]/g, "![$1]($1)");
    // [[foo.extension]]
    source = source.replace(/\[\[([^|\]]+)\.([^|\]\.]+)\]\]/g, "[$1]($1.$2)");
    // [[foo]]
    source = source.replace(/\[\[([^|\]]+)\]\]/g, "[$1]($1.md)");
    // [[foo|bar.extension]]
    source = source.replace(/\[\[([^|\]]+)\|([^\]]+)\.([^\]\.]+)\]\]/g, "[$1]($2.$3)");
    source = source.replace(/\[\[([^|\]]+)\|([^\]]+)\]\]/g, "[$1]($2.md)");


    var compiled = marked(source, {renderer: this.markdownRenderer});

    // Replace equations
    for (var mEq in multilineEqs) {
      compiled = compiled.replace(mEq, '<div><script type="math/tex;mode=display">' + multilineEqs[mEq] + '</script></div>');
    }
    for (var iEq in inlineEqs) {
      compiled = compiled.replace(iEq, '<span><script type="math/tex">' + inlineEqs[iEq] + '</script></span>');
    }

    return compiled;
  }),
});
