import { formatError } from './validation';

const justMsg = { msg: 'Just an error message' };
test('just a message', () => {
  expect(formatError('source.js', justMsg)).toEqual(`source.js: ${justMsg.msg}`);
});

const msgAndPath = { msg: 'This field was wrong', path: '<root>.foo.bar' };
test('a message with a path', () => {
  expect(formatError('source.js', msgAndPath)).toEqual(`source.js: ${msgAndPath.msg} at ${msgAndPath.path}`);
});

const msgStart = {
  msg: 'The field here was wrong',
  start: { line: 3, column: 20 },
};
test('message with start location', () => {
  expect(formatError('source.js', msgStart)).toEqual(`source.js:3.20: ${msgStart.msg}`);
});

const all = {
  msg: 'The field here was wrong',
  path: '<root>.foo.bar',
  start: { line: 3, column: 20 },
  end: { line: 5, column: 0 },
};
test('', () => {
  expect(formatError('source.js', all)).toEqual(`source.js:3.20-5.0: ${all.msg} at ${all.path}`);
});
