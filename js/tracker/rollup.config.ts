import resolve from 'rollup-plugin-node-resolve'
import commonjs from 'rollup-plugin-commonjs'
import uglify from 'rollup-plugin-uglify';
import {minify} from 'uglify-es';

const pkg = require('./package.json')

const libraryName = 'show-ad'

export default [{
  input: `compiled/index.js`,
  output: [
    {file: pkg.main, format: 'es'},
  ],
  sourcemap: false,
  // Indicate here external modules you don't wanna include in your bundle (i.e.: 'lodash')
  external: [],
  watch: {
    include: 'compiled/**',
  },
  plugins: [
    // Allow bundling cjs modules (unlike webpack, rollup doesn't understand cjs)
    commonjs(),
    // Allow node_modules resolution, so you can use 'external' to control
    // which external modules to include in the bundle
    // https://github.com/rollup/rollup-plugin-node-resolve#usage
    resolve(),
    uglify({}, minify)
  ],
}]
