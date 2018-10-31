const Container = function(name, image) {
  return {
    name,
    image,
  }
};

const Deployment = function(name, replicas, containers) {
  return {
    apiVersion: 'apps/v1',
    kind: 'Deployment',
    metadata: {
      name,
      labels: {
	app: name,
      },
    },
    spec: {
      replicas,
      selector: {
	matchLabels: {
	  app: name,
	},
      },
      template: {
	metadata: {
	  labels: {
	    app: name,
	  },
	},
	containers,
      },
    },
  };
};


export default {
  Container,
  Deployment,
};
