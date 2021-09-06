
const { sep, join, normalize, resolve } = require('path')
const fs = require('fs')

function getAppDir() {
    let dir = process.cwd()

    while (dir.length && dir[dir.length - 1] !== sep) {
        if (fs.existsSync(join(dir, 'package.json'))) {
            return dir
        }

        dir = normalize(join(dir, '..'))
    }
}

const appDir = join(getAppDir(), 'frontend'),
    store = join(appDir, 'store'),
    cliDir = resolve(__dirname, '..')

module.exports = {
    cli: cliDir,
    store: store,
    resolve: {
        cli: dir => join(cliDir, dir),
        store: dir => join(store, dir)
    }
}