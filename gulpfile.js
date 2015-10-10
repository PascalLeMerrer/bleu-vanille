var gulp    = require('gulp');
var cssnext = require("gulp-cssnext");
var watch   = require('gulp-watch');
var cucumber = require('gulp-cucumber');

gulp.task('styles', function() {
  gulp.src("./src/bleuvanille/public/css/marketing.css")
  .pipe(watch("./src/bleuvanille/public/css/marketing.css"))
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


gulp.task('test', function() {
    return gulp.src('tests/*')
			.pipe(cucumber({
				'steps': 'tests/step_definitions/*.js',
				'format': 'pretty'
			}));
});

/*
 * Build assets by default.
 */
gulp.task('default', ['styles', 'images', 'js', 'fonts']);
