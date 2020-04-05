# blog

## プロジェクトの初期化

```
go mod init $(basename $(pwd))
```

## gulp

```
mkdir web
cd web
npm i -g gulp-cli
npm init -y
npm i -D gulp @babel/core @babel/register gulp-babel gulp-uglify gulp-rename gulp-livereload gulp-less gulp-postcss autoprefixer autorem cssnano gulp-util del
```
