import handlebars from 'handlebars/lib/handlebars';

const source = `
<div class="entry">
  <h1>{{title}}</h1>
  <div class="body">
    {{body}}
  </div>
</div>
`;

const template = handlebars.compile(source);
const context = { title: 'My New Post', body: 'This is my first post!' };

export default [
  { file: 'index.html', value: template(context) },
];
