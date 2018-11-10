import includePaths from 'rollup-plugin-includepaths';

const includePathOptions = {
    include: {},
    paths: ['.'],
    external: [],
    extensions: ['.js']
};

export default {
    input: './std.js',
    output: {
        format: 'es',
    },
    plugins: [ includePaths(includePathOptions) ],
};
