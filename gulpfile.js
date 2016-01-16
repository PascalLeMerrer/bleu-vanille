var gulp     = require('gulp');
var cssnext  = require("gulp-cssnext");
var watch    = require('gulp-watch');
var riot     = require('gulp-riot');
var cucumber = require('gulp-cucumber');
var argv     = require('yargs').argv;

gulp.task('styles', function() {
  gulp.src("./src/bleuvanille/public/css/*.css")
  .pipe(watch("./src/bleuvanille/public/css/*.css"))
  .pipe(cssnext({
    compress: false,  // default is false
  }))
  .pipe(gulp.dest("./public/css"))
});

gulp.task('images', function() {
  gulp.src('./src/bleuvanille/public/img/**/*.*', {overwrite: false})
  .pipe(watch('./src/bleuvanille/public/img/**/*.*'))
  .pipe(gulp.dest('./public/img'));
});

gulp.task('fonts', function() {
  gulp.src('./src/bleuvanille/public/fonts/**/*.*', {overwrite: false})
  .pipe(watch('./src/bleuvanille/public/fonts/**/*.*'))
  .pipe(gulp.dest('./public/fonts'));
});

gulp.task('js', function() {
  gulp.src('./src/bleuvanille/public/js/**/*.js', {overwrite: true})
    .pipe(watch('./src/bleuvanille/public/js/**/*.js'))
    .pipe(gulp.dest('./public/js'));
});

gulp.task('html', function() {
  gulp.src('./src/bleuvanille/public/html/**/*.html', {overwrite: true})
    .pipe(watch('./src/bleuvanille/public/html/**/*.html'))
    .pipe(gulp.dest('./public/html'));
});

gulp.task('publicroot', function() {
  gulp.src('./src/bleuvanille/public/*', {overwrite: true})
    .pipe(watch('./src/bleuvanille/public/*'))
    .pipe(gulp.dest('./public/'));
});

gulp.task('riot', function() {
  gulp.src('./src/bleuvanille/public/tags/**/*.html', {overwrite: true})
    .pipe(watch('./src/bleuvanille/public/tags/**/*.html'))
    .pipe(riot())
    .pipe(gulp.dest('./public/tags'));
});

gulp.task('test', function() {
  var tags = '';
  var parameters;
  if(argv.tags) {
    parameters = argv.tags.split(' ')
  } else if (argv.t) {
    parameters = argv.t.split(' ')
  }

  if(parameters) {
    for (var i = 0; i < parameters.length; i++) {
      if (i > 0) {
        tags += ','
      }
      tags += '@' + parameters[i]
    }
  }
  
  return gulp.src('tests/*')
		.pipe(cucumber({
			'steps': 'tests/step_definitions/*.js',
			'format': 'pretty'
      ,tags: tags
		}));
});

gulp.task('dist', function() {
  gulp.src("./src/bleuvanille/public/css/*.css")
  .pipe(cssnext({
    compress: true
  }))
  .pipe(gulp.dest("./dist/public/css"))

  gulp.src('./src/bleuvanille/public/img/**/*.*', {overwrite: true})
  .pipe(gulp.dest('./dist/public/img'));

  gulp.src('./src/bleuvanille/public/fonts/**/*.*', {overwrite: true})
  .pipe(gulp.dest('./dist/public/fonts'));

  gulp.src('./src/bleuvanille/public/js/**/*.js', {overwrite: true})
    .pipe(gulp.dest('./dist/public/js'));

  gulp.src('./src/bleuvanille/public/html/**/*.html', {overwrite: true})
    .pipe(gulp.dest('./dist/public/html'));

  gulp.src('./src/bleuvanille/public/tags/**/*.html', {overwrite: true})
    .pipe(riot())
    .pipe(gulp.dest('./dist/public/tags'));
});

/*
 * Build assets by default.
 */
gulp.task('default', ['riot', 'html', 'styles', 'images', 'js', 'fonts', 'publicroot']);
