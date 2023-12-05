import { expect, afterAll, beforeAll, test } from 'bun:test'
// @ts-ignore
import { createNetwork, type Network } from '@eth-optimism/mocktimism'
import { createClient, http } from 'viem'

let network: Network

beforeAll(async () => {
	network = createNetwork({
		chains: [
			{
				name: 'mainnet',
				forkUrl: 'https://mainnet.infura.io/v3/420',
				chainId: 1,
			},
			{
				name: 'optimism',
				forkUrl: 'https://mainnet.optimism.io',
				chainId: 10,
				baseChainId: 1,
			}
		]
	})
	await network.start()
})

afterAll(async () => {
	await network.stop()
})

test('should be able to create viem clients', async () => {
	const l1Client = createClient({
		transport: http(network.getChain('mainnet').rpcUrl),
		chain: network.getChain('mainnet'),
	})
	expect(await l1Client.request({
		method: 'eth_chainId'
	})).toBe(network.getChain('mainnet').id)
}) 
