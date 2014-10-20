var gulp   = require('gulp');
var util   = require('gulp-util');
var size = require('gulp-size');
var concat = require('gulp-concat');
var jshint = require('gulp-jshint');
var sass = require('gulp-sass');
var connect = require('gulp-connect');
var watch = require('gulp-watch');
var traceur = require('gulp-traceur');
var uglify = require('gulp-uglify');
var minifyCss = require('gulp-minify-css');
var url = require('url');
var proxy = require('proxy-middleware');

var config = {
  production: process.env.ENV === 'production',

  js: {
    src: [
      'js/models/page.js',
      'js/models/folder.js',
      'js/views/folder_item_view.js',
      'js/views/page_item_view.js',
      'js/views/folder_list_view.js',
      'js/views/page_view.js',
      'js/app.js',
    ],
    dest: 'dist/assets',
  },

  css: {
    src: 'style/*.scss',
    dest: 'dist/assets',
  },

  html: {
    src: 'index.html',
    dest: 'dist',
  },

  imgs: {
    src: 'images/*',
    dest: 'dist/assets',
  },

  vendor_js: {
    src: [
      'bower_components/jquery/dist/jquery.js',
      'bower_components/underscore/underscore.js',
      'bower_components/backbone/backbone.js',
      'bower_components/marked/lib/marked.js',
      'bower_components/highlightjs/highlight.pack.js',
    ],
    dest: 'dist/assets',
  },

  vendor_css: {
    src: [
      'bower_components/highlightjs/styles/tomorrow.css',
    ],
    dest: 'dist/assets',
  },

  fonts: {
    src: 'bower_components/fontawesome/fonts/fontawesome-webfont.woff',
    dest: 'dist/assets/fonts'
  }
};

gulp.task('js', function () {
  return gulp.src(config.js.src)
    .pipe(traceur({
      modules: 'inline',
      moduleName: function (path) {
        return path;
      }
    }))
    .pipe(concat({path: 'app.js'}))
    .pipe(config.production ? uglify() : util.noop())
    .pipe(gulp.dest(config.js.dest))
    .pipe(connect.reload());
});

gulp.task('css', function () {
  return gulp.src(config.css.src)
    .pipe(concat({path: 'app.css'}))
    .pipe(sass())
    .pipe(config.production ? minifyCss() : util.noop())
    .pipe(gulp.dest(config.css.dest))
    .pipe(connect.reload());
});

gulp.task('lint', function() {
  return gulp.src(config.js.src)
    .pipe(jshint())
    .pipe(jshint.reporter('default'));
});

gulp.task('html', function () {
  return gulp.src(config.html.src)
    .pipe(gulp.dest(config.html.dest))
    .pipe(connect.reload());
});

gulp.task('imgs', function () {
  return gulp.src(config.imgs.src)
    .pipe(gulp.dest(config.imgs.dest));
});

gulp.task('fonts', function () {
  return gulp.src(config.fonts.src)
    .pipe(gulp.dest(config.fonts.dest));
});

gulp.task('vendor-js', function () {
  return gulp.src(config.vendor_js.src)
    .pipe(concat({path: 'vendor.js'}))
    .pipe(config.production ? uglify() : util.noop())
    .pipe(gulp.dest(config.vendor_js.dest));
});

gulp.task('vendor-css', function () {
  return gulp.src(config.vendor_css.src)
    .pipe(concat({path: 'vendor.css'}))
    .pipe(config.production ? minifyCss() : util.noop())
    .pipe(gulp.dest(config.vendor_css.dest));
});

gulp.task('server', function() {
  connect.server({
    root: 'dist',
    fallback: 'index.html',
    livereload: true,
    middleware: function(connect, o) {
      var options = url.parse('http://localhost:3000/v1');
      options.route = '/v1';
      return [proxy(options)];
    }
  });
});

gulp.task('watch', ['default', 'server'], function () {
  gulp.watch([config.js.src, config.css.src, config.html.src], ['html', 'js', 'css']);
});


gulp.task('default', ['lint', 'js', 'css', 'html', 'vendor-js', 'vendor-css', 'imgs', 'fonts'], function () {
  return gulp.src("dist/**")
    .pipe(size({
      showFiles: true
    }));
});
