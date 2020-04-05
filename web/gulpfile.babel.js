const { src, dest, watch } = require('gulp');
const babel = require('gulp-babel');
const uglify = require('gulp-uglify');
const rename = require('gulp-rename');
const livereload = require('gulp-livereload');

function javascript(cb) {
  src('src/**/*.js')
    .pipe(babel())
    .pipe(src('vendor/*.js'))
    .pipe(dest('dist/'))
    .pipe(uglify())
    .pipe(rename({ extname: '.min.js' }))
    .pipe(dest('dist/'));
  cb();
}

function css(cb) {
  src('src/**/*.less')
    .pipe(
      less().on('error', function (err) {
        util.log(err);
        this.emit('end');
      })
    )
    .pipe(
      postcss([
        autoprefixer({ overrideBrowserslist: ['last 4 versions'] }),
        cssnano({ preferredQuote: 'single' }),
        autorem(),
      ])
    )
    .pipe(rename({ extname: '.min.css' }))
    .pipe(dest('dist/'));
  cb();
}

exports.default = function () {
  livereload.listen(8080, 'localhost');
  watch('src/**/*.js', javascript);
  watch('src/**/*.less', css);
};

exports.watch = function () {
  livereload.listen(8080, 'localhost');
};
