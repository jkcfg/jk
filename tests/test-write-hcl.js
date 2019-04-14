import std from '@jkcfg/std';

// Github organization
const organization = 'myorg';

const config = {
  provider: {
    github: {
      organization,
    },
  },
  github_membership: {
    myorg_foo: {
      username: 'foo',
      role: 'admin',
    },
  },
};

std.write(config, 'github.tf');
