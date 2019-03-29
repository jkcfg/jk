import { mix, patch, merge } from '../std/std_merge';

test('mix objects', () => {
  const r = mix({foo: 1}, {bar: 2}, {foo: 3});

  expect(r).toEqual({
    foo: 3,
    bar: 2
  });
});

test('mix transforms', () => {
  const addLabel = (k, v) => {
    const labels = {};
    labels[k] = v;
    return o => patch(o, { labels });
  };
  const orig = { foo: 1, labels: { l1: 'v1', l2: 'v2' } };

  expect(mix(orig, addLabel('l3', 'v3'), addLabel('l1', 'w1'))).toEqual({
    foo: 1,
    labels: {
      l1: 'w1',
      l2: 'v2',
      l3: 'v3',
    }
  });

  // orig has been left untouched.
  expect(orig).toEqual(
    { foo: 1, labels: { l1: 'v1', l2: 'v2' } }
  );
});

test('trivial patch', () => {
  expect(patch({}, {})).toEqual({});
  expect(patch({}, {foo: 1})).toEqual({foo: 1});
  expect(patch({}, {foo: {bar: "baz"}})).toEqual({foo: {bar: "baz"}});
  expect(patch({foo: 1}, {foo: {bar: "baz"}})).toEqual({foo: {bar: "baz"}});
});

test('nested patch', () => {
  const orig = {
    foo: {bar: 1},
    baz: 2
  };

  expect(patch(orig, {foo: {bar: 3}})).toEqual(
    {
      foo: {bar: 3},
      baz: 2
    }
  );

  expect(patch(orig, {foo: {bar: 3, baz: 4}})).toEqual(
    {
      foo: {bar: 3, baz: 4},
      baz: 2
    }
  );

  // orig has been left untouched.
  expect(orig).toEqual({
    foo: {bar: 1},
    baz: 2
  });
});

test('merge values', () => {
  expect(merge(1, 2)).toEqual(3);
  expect(merge('a', 'b')).toEqual('ab');
  expect(merge('a', 1)).toEqual('a1');
  expect(merge(2, 'b')).toEqual('2b');
  expect(merge('s', true)).toEqual('strue');
  expect(merge('str ', {foo: 'bar'})).toEqual(`str {"foo":"bar"}`);
  expect(merge('str ', [2, 3, true])).toEqual(`str [2,3,true]`);
  expect(merge([1,2], [3,4])).toEqual([1,2,3,4]);
  expect(merge({foo: 1}, {bar: 2})).toEqual({foo: 1, bar: 2});
});

test('deep merge', () => {
  const orig = {
    foo: {
      bar: 'bar',
      boo: 'boo',
    },
    baz: 'baz',
  };

  // regular assign syntax (no '+')
  expect(merge(orig, {
    foo: {bar: 'replaced'},
  })).toEqual({
    foo: {bar: 'replaced'},
    baz: 'baz',
  });

  // deep merge syntax ('+')
  expect(merge(orig, {
    'foo+': {bar: 'replaced'},
  })).toEqual({
    foo: {
      bar: 'replaced',
      boo: 'boo',
    },
    baz: 'baz',
  });

  // concat values
  expect(merge(orig, {
    'foo+': { 'bar+': 'concat' },
  })).toEqual({
    foo: {
      bar: 'barconcat',
      boo: 'boo',
    },
    baz: 'baz'
  });

  // original untouched
  expect(orig).toEqual({
    foo: {
      bar: 'bar',
      boo: 'boo',
    },
    baz: 'baz',
  });
});
