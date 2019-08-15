import {
  mergeFull, deep, first, replace, deepWithKey,
} from './merge';

test('mergeFull: default merging of primitive values', () => {
  expect(mergeFull(1, 2)).toEqual(2);
  expect(mergeFull('a', 'b')).toEqual('b');
  expect(() => mergeFull('a', 1)).toThrow();
  expect(() => mergeFull(true, 'b')).toThrow();
  expect(mergeFull([1, 2], [3, 4])).toEqual([3, 4]);
  expect(mergeFull({ foo: 1 }, { bar: 2 })).toEqual({ foo: 1, bar: 2 });
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

test('mergeFull: array of objects, merging objects identified by a key', () => {
  const result = mergeFull(pod, sidecarImage, {
    spec: deep({
      containers: deepWithKey('name'),
    }),
  });

  expect(result.spec.containers.length).toEqual(2);
  expect(result.spec.containers[1].image).toEqual('sidecar:v2');
});

test('mergeFull: pick the deep merge strategy when encountering an object as rule', () => {
  const result = mergeFull(pod, sidecarImage, {
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

  expect(() => mergeFull(pod, sidecarImageNotObject, rules)).toThrow();
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

  expect(() => mergeFull(pod, sidecarImageNotArray, rules)).toThrow();
});

test('first: basic', () => {
  const result = mergeFull(pod, sidecarImage, {
    spec: {
      containers: first(),
    },
  });

  expect(result.spec.containers.length).toEqual(2);
  expect(result.spec.containers[0].name).toEqual('my-app');
  expect(result.spec.containers[1].image).toEqual('sidecar:v1');
});

test('replace: basic', () => {
  const result = mergeFull(pod, sidecarImage, {
    spec: {
      containers: replace(),
    },
  });

  expect(result.spec.containers.length).toEqual(1);
  expect(result.spec.containers[0].name).toEqual('sidecar');
});
