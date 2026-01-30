import {cpSync, mkdirSync} from 'fs';

const externals: Record<string, string> = {
    'react': 'React',
    'redux': 'Redux',
    'react-dom': 'ReactDOM',
    'react-redux': 'ReactRedux',
    'prop-types': 'PropTypes',
    'react-bootstrap': 'ReactBootstrap',
};

const externalPattern = Object.keys(externals).map((k) => k.replace(/[-/]/g, '\\$&')).join('|');
const externalFilter = new RegExp(`^(${externalPattern})$`);

const result = await Bun.build({
    entrypoints: ['./src/index.jsx'],
    outdir: './dist',
    naming: 'main.js',
    sourcemap: 'external',
    target: 'browser',
    format: 'iife',
    plugins: [
        {
            name: 'mattermost-externals',
            setup(build) {
                build.onResolve({filter: externalFilter}, (args) => ({
                    path: args.path,
                    namespace: 'mattermost-external',
                }));

                build.onLoad({filter: /.*/, namespace: 'mattermost-external'}, (args) => ({
                    contents: `module.exports = window["${externals[args.path]}"];`,
                    loader: 'js',
                }));
            },
        },
        {
            name: 'css-loader',
            setup(build) {
                build.onLoad({filter: /\.css$/}, async (args) => {
                    const css = await Bun.file(args.path).text();
                    return {
                        contents: [
                            `const css = ${JSON.stringify(css)};`,
                            'if (typeof document !== "undefined") {',
                            '    const style = document.createElement("style");',
                            '    style.textContent = css;',
                            '    document.head.appendChild(style);',
                            '}',
                            'export default css;',
                        ].join('\n'),
                        loader: 'js',
                    };
                });
            },
        },
        {
            name: 'svg-inline-loader',
            setup(build) {
                build.onLoad({filter: /\.svg$/}, async (args) => {
                    const svg = await Bun.file(args.path).text();
                    return {
                        contents: `export default ${JSON.stringify(svg)};`,
                        loader: 'js',
                    };
                });
            },
        },
    ],
});

if (!result.success) {
    console.error('Build failed:');
    for (const log of result.logs) {
        console.error(log);
    }
    process.exit(1);
}

mkdirSync('./dist/static', {recursive: true});
cpSync('./src/assets/images', './dist/static', {recursive: true});

console.log('Build complete');
