export default function validate(obj) {
  if (obj.name === 'Valid') return true;
  return 'object name is not "Valid"';
}
