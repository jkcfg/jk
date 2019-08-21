import {
  merge, deep, first, replace, deepWithKey,
} from './merge';

test('merge: default merging of primitive values', () => {
  expect(merge(1, 2)).toEqual(2);
  expect(merge('a', 'b')).toEqual('b');
  expect(() => merge('a', 1)).toThrow();
  expect(() => merge(true, 'b')).toThrow();
  expect(merge([1, 2], [3, 4])).toEqual([3, 4]);
  expect(merge({ foo: 1 }, { bar: 2 })).toEqual({ foo: 1, bar: 2 });
});

const pod = {
  spec: {
    containers: [{
      name: 'my-app',
      image: 'busybox',
      command: ['sh', '-c', 'echo Hello Kubernetes! && sleep 3600'],
    }, {
      name: 'sidecar',
      image: 'sidecar:v1',
    }],
  },
};

const sidecarImage = {
  spec: {
    containers: [{
      name: 'sidecar',
      image: 'sidecar:v2',
    }],
  },
};

test('merge: array of objects, merging objects identified by a key', () => {
  const result = merge(pod, sidecarImage, {
    spec: deep({
      containers: deepWithKey('name'),
    }),
  });

  expect(result.spec.containers.length).toEqual(2);
  expect(result.spec.containers[1].image).toEqual('sidecar:v2');
});

test('merge: pick the deep merge strategy when encountering an object as rule', () => {
  const result = merge(pod, sidecarImage, {
    spec: {
      containers: deepWithKey('name'),
    },
  });

  expect(result.spec.containers.length).toEqual(2);
  expect(result.spec.containers[1].image).toEqual('sidecar:v2');
});

test('deep: throw on wrong input type', () => {
  const sidecarImageNotObject = {
    spec: [{
      containers: [{
        name: 'sidecar',
        image: 'sidecar:v2',
      }],
    }],
  };

  const rules = {
    spec: deep({
      containers: deepWithKey('name'),
    }),
  };

  expect(() => merge(pod, sidecarImageNotObject, rules)).toThrow();
});

test('deepWithKey: throw on wrong input type', () => {
  const sidecarImageNotArray = {
    spec: {
      containers: {
        name: 'sidecar',
        image: 'sidecar:v2',
      },
    },
  };

  const rules = {
    spec: deep({
      containers: deepWithKey('name'),
    }),
  };

  expect(() => merge(pod, sidecarImageNotArray, rules)).toThrow();
});

test('first: basic', () => {
  const result = merge(pod, sidecarImage, {
    spec: {
      containers: first(),
    },
  });

  expect(result.spec.containers.length).toEqual(2);
  expect(result.spec.containers[0].name).toEqual('my-app');
  expect(result.spec.containers[1].image).toEqual('sidecar:v1');
});

test('replace: basic', () => {
  const result = merge(pod, sidecarImage, {
    spec: {
      containers: replace(),
    },
  });

  expect(result.spec.containers.length).toEqual(1);
  expect(result.spec.containers[0].name).toEqual('sidecar');
});
