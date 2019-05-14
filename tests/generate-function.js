export default function f(message = 'success') {
  return [
    { file: 'object.yaml', value: { message } },
  ];
}
