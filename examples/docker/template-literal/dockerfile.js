import { Dockerfile } from './docker';

const myService = {
  name: 'my-service',
  port: 80,
};

export default [
  { path: 'Dockerfile', value: Dockerfile(myService) },
];
