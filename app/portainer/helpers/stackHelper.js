import _ from 'lodash-es';
import YAML from 'yaml';
import GenericHelper from '@/portainer/helpers/genericHelper';
import { ExternalStackViewModel } from '@/react/docker/stacks/view-models/external-stack';

angular.module('portainer.app').factory('StackHelper', [
  function StackHelperFactory() {
    'use strict';
    var helper = {};

    helper.getExternalStacksFromContainers = function (containers) {
      return getExternalStacksFromLabel(containers, 'com.docker.compose.project', 2);
    };

    helper.getExternalStacksFromServices = function (services) {
      return getExternalStacksFromLabel(services, 'com.docker.stack.namespace', 1);
    };

    function getExternalStacksFromLabel(items, label, type) {
      return _.uniqBy(
        items.filter((item) => item.Labels && item.Labels[label]).map((item) => new ExternalStackViewModel(item.Labels[label], type, item.Created)),
        'Name'
      );
    }

    helper.validateYAML = validateYAML;

    return helper;
  },
]);

function validateYAML(yaml, containerNames, originalContainersNames = []) {
  let yamlObject;

  try {
    yamlObject = YAML.parse(yaml, { mapAsMap: true, maxAliasCount: 10000 });
  } catch (err) {
    return 'There is an error in the yaml syntax: ' + err;
  }

  const names = _.uniq(GenericHelper.findDeepAll(yamlObject, 'container_name'));

  const duplicateContainers = _.intersection(_.difference(containerNames, originalContainersNames), names);

  if (duplicateContainers.length === 0) {
    return '';
  }

  return (
    (duplicateContainers.length === 1 ? 'This container name is' : 'These container names are') +
    ' already used by another container running in this environment: ' +
    _.join(duplicateContainers, ', ') +
    '.'
  );
}

export function extractContainerNames(yaml = '') {
  let yamlObject;

  try {
    yamlObject = YAML.parse(yaml, { maxAliasCount: 10000 });
  } catch (err) {
    return [];
  }

  return _.uniq(GenericHelper.findDeepAll(yamlObject, 'container_name'));
}
