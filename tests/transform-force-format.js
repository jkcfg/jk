// Reminder: each file will be treated as a stream, and the function
// called on each value.

export default function (x) {
  if (Array.isArray(x)) {
    x.push({ seen: true });
  } else {
    x.seen = true;
  }
}
