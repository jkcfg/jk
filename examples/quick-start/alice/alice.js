// Define a developer.
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

// Write the developer description as YAML.
export default [
  { value: alice, file: `developers/${alice.name.toLowerCase()}.yaml` },
];
