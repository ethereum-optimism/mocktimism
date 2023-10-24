import { ExampleContract } from './ExampleContract.sol'
import { Address, createPublicClient, createWalletClient, http } from 'viem'
import { mainnet, optimism, optimismGoerli } from 'viem/chains'
import { mnemonicToAccount } from 'viem/accounts'
import { publicL1OpStackActions, publicL2OpStackActions, walletL1OpStackActions, walletL2OpStackActions } from 'op-viem'

const account = mnemonicToAccount('test test test test test test test test test test test junk')

const clients = {
	public: {
		[mainnet.id]: createPublicClient({
			chain: mainnet,
			transport: http('http://localhost:8545'),
		}).extend(publicL1OpStackActions),
		[optimism.id]: createPublicClient({
			chain: optimism,
			transport: http('http://localhost:9545'),
		}).extend(publicL2OpStackActions)
	},
	wallet: {
		[mainnet.id]: createWalletClient({
			chain: mainnet,
			transport: http('http://localhost:8545'),
			account
		}).extend(walletL1OpStackActions),
		[optimism.id]: createWalletClient({
			chain: optimism,
			transport: http('http://localhost:9545'),
			account
		}).extend(walletL2OpStackActions)
	}
}

export const mintOnL2 = async () => {
	const tokenId = BigInt(420420)

	console.info(`minting tokeId ${tokenId.toString()} on l2...`)
	// TODO update version of op viem to mint on l2 with writeContractDeposit
	const l1TxHash = await clients.wallet[mainnet.id].writeContract({
		abi: ExampleContract.abi,
		address: '0x1df10ec981ac5871240be4a94f250dd238b77901',
		functionName: 'mint',
		args: [tokenId]
	}).catch(e => {
		console.error(e)
		throw new Error('l1: writeContractDeposit failed')
	})
	console.info('l1TxHash', l1TxHash)

	console.info('waiting for l1 receipt...')
	// Wait for l1 to confirm
	const l1TxReceipt = await clients.public[mainnet.id].waitForTransactionReceipt({ hash: l1TxHash }).catch(e => {
		console.error(e)
		throw new Error('l1: waitForTransactionReceipt failed')
	})
	console.info('l1TxReceipt', l1TxReceipt)

	console.info('getting l2 tx hash...')
	// get the deterministic l2 tx hash
	const l2TxHashes = await clients.public[mainnet.id].getL2HashesForDepositTx({ l1TxReceipt }).catch(e => {
		console.error(e)
		throw new Error('l1: getL2HashesForDepositTx failed')
	})
	console.info('l2TxHashes', l2TxHashes)

	console.info('waiting for l2 receipt...')
	// wait for l2 to confirm now
	const l2TxReceipt = await clients.public[optimism.id].waitForTransactionReceipt({ hash: l2TxHashes[0] }).catch(e => {
		console.error(e)
		throw new Error('l2: waitForTransactionReceipt failed')
	})
	console.info('l2TxReceipt', l2TxReceipt)

	console.info('confirming ownerof')
	await clients.public[optimism.id].readContract({
		abi: ExampleContract.abi,

		address: '0x1df10ec981ac5871240be4a94f250dd238b77901',
		functionName: 'ownerOf',
		args: [tokenId]
	}).then(owner => {
		if (owner !== account.address) {
			throw new Error(`ownerOf doesn't match expected: ${account.address} received: ${owner}`)
		}
		return owner
	}).catch(e => {
		console.error(e)
		throw new Error('ownerOf threw error')
	})

	return l2TxReceipt
}

if (import.meta.main) {
	mintOnL2().then(console.log).catch(e => {
		console.error(e)
	})
}
