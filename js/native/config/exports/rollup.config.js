// rollup.config.js
import commonjs from 'rollup-plugin-commonjs';
import nodeResolve from 'rollup-plugin-node-resolve';
import alias from 'rollup-plugin-alias';
import replace from 'rollup-plugin-replace';
import fs from "fs";

const style = fs.readFileSync("./build/style.css");


const substituteModulePaths = {
    'crypto': 'build/module/adapters/crypto.browser.js',
    'hash.js': 'build/temp/hash.js'
}

export default {
    entry: 'build/module/index.js',
    sourceMap: true,
    plugins: [
        alias(substituteModulePaths),
        nodeResolve({
            browser: true
        }),
        commonjs(),
	    replace({
		    __STYLE_TEMPLATE__: style.toString().replace(/"/ig,`\\"`),
	    })
    ]
}
