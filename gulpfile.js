var gulp = require('gulp');
var path = require('path');
var shell = require('gulp-shell');

var goPath = "./mypackage/**/*.go"

gulp.task("compilepkg", function() {
    return gulp.src(goPath, {
            read: false
        })
        .pipe(shell(['go install .'], {
            templateData: {
                stripPath: function(filePath) {
                    var subPath = filePath.substring(process.cwd().length + 1);

                    var pkg = subPath.substring(0, subPath.lastIndexOf(path.sep));
                    return pkg
                }
            }
        }));
});


gulp.task('watch', function() {
    gulp.watch(goPath, ['compilepkg']);

})
