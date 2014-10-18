var gulp   = require('gulp');
var concat = require('gulp-concat');
var jshint = require('gulp-jshint');
var sass = require('gulp-sass');
var connect = require('gulp-connect');
var watch = require('gulp-watch');
var traceur = require('gulp-traceur');

var config = {
  js: {
    src: [
      'js/models/*.js',
      'js/views/*.js',
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
    src: 'logo.svg',
    dest: 'dist/assets',
  },

  vendor_js: {
    src: [
      'bower_components/jquery/dist/jquery.js',
      'bower_components/underscore/underscore.js',
      'bower_components/backbone/backbone.js',
      'bower_components/marked/lib/marked.js',
    ],
    dest: 'dist/assets',
  },
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
    .pipe(gulp.dest(config.js.dest))
    .pipe(connect.reload());
});

gulp.task('css', function () {
  return gulp.src(config.css.src)
    .pipe(concat({path: 'app.css'}))
    .pipe(sass())
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
    .pipe(gulp.dest(config.imgs.dest))
    .pipe(connect.reload());
});

gulp.task('vendor-js', function () {
  return gulp.src(config.vendor_js.src)
    .pipe(concat({path: 'vendor.js'}))
    .pipe(gulp.dest(config.vendor_js.dest));
});

gulp.task('server', function() {
  connect.server({
    root: 'dist',
    fallback: 'index.html',
    livereload: true,
  });
});

gulp.task('watch', ['server'], function () {
  gulp.watch([config.js.src, config.css.src, config.html.src], ['html', 'js', 'css', 'imgs']);
});


gulp.task('default', ['lint', 'js', 'css', 'html', 'vendor-js']);
