export default function f(message = 'success') {
  return [
    { path: 'object.yaml', value: { message } },
  ];
}
