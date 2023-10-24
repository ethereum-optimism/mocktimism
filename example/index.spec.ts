import { ChildProcess, exec } from 'child_process'
import { mintOnL2 } from './index.js'
import { afterAll, beforeAll, test } from 'bun:test'
import { join } from 'path'
import waitOn from 'wait-on'

let mocktimismProcess: ChildProcess

beforeAll(async () => {
	mocktimismProcess = exec(`go run ../cmd/main.go up --config ${join(__dirname, 'mocktimism.toml')}`, (err, stdout, stderr) => {
		if (err) {
			console.error(err)
			throw new Error('failed to start mocktimism')
		}
		console.info(stdout)
		console.error(stderr)
	})

	await waitOn({
		resources: [
			'http://localhost:8545',
			'http://localhost:9545',
		]
	})
})

afterAll(async () => {
	mocktimismProcess.kill()
})

test(mintOnL2.name, async () => {
	await mintOnL2()
}) 
