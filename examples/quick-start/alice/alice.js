// Alice is a developer.
const alice = {
  name: 'Alice',
  beverage: 'Club-Mate',
  monitors: 2,
  languages: [
    'python',
    'haskell',
    'c++',
    '68k assembly', // Alice is cool like that!
  ],
};

// Instruct to write the alice object as a YAML file.
export default [
  { value: alice, path: `developers/${alice.name.toLowerCase()}.yaml` },
];
