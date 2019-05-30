import { Dockerfile } from './docker';

const myService = {
  name: 'my-service',
  port: 80,
};

export default [
  { file: 'Dockerfile', value: Dockerfile(myService) },
];
