import { ContainerId } from '@/react/docker/containers/types';
import axios, {
  agentTargetHeader,
  parseAxiosError,
} from '@/portainer/services/axios';
import { EnvironmentId } from '@/react/portainer/environments/types';

import { NetworkId, DockerNetwork } from './types';

type NetworkAction = 'connect' | 'disconnect' | 'create';

export async function getNetwork(
  environmentId: EnvironmentId,
  networkId: NetworkId,
  { nodeName }: { nodeName?: string } = {}
) {
  try {
    const { data: network } = await axios.get<DockerNetwork>(
      buildUrl(environmentId, networkId),
      nodeName
        ? {
            headers: {
              [agentTargetHeader]: nodeName,
            },
          }
        : undefined
    );
    return network;
  } catch (e) {
    throw parseAxiosError(e as Error, 'Unable to retrieve network details');
  }
}

export async function deleteNetwork(
  environmentId: EnvironmentId,
  networkId: NetworkId,
  { nodeName }: { nodeName?: string } = {}
) {
  try {
    await axios.delete(
      buildUrl(environmentId, networkId),
      nodeName
        ? {
            headers: {
              [agentTargetHeader]: nodeName,
            },
          }
        : undefined
    );
    return networkId;
  } catch (e) {
    throw parseAxiosError(e as Error, 'Unable to remove network');
  }
}

export async function disconnectContainer(
  environmentId: EnvironmentId,
  networkId: NetworkId,
  containerId: ContainerId,
  nodeName?: string
) {
  try {
    await axios.post(
      buildUrl(environmentId, networkId, 'disconnect'),
      {
        Container: containerId,
        Force: false,
      },
      nodeName
        ? {
            headers: {
              [agentTargetHeader]: nodeName,
            },
          }
        : undefined
    );
    return { networkId, environmentId };
  } catch (e) {
    throw parseAxiosError(
      e as Error,
      'Unable to disconnect container from network'
    );
  }
}

function buildUrl(
  environmentId: EnvironmentId,
  networkId?: NetworkId,
  action?: NetworkAction
) {
  let url = `endpoints/${environmentId}/docker/networks`;

  if (networkId) {
    url += `/${networkId}`;
  }

  if (action) {
    url += `/${action}`;
  }

  return url;
}
