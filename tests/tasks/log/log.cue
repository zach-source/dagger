package main

import (
	"dagger.io/dagger"
	"dagger.io/dagger/core"
)

dagger.#Plan & {
	client: logger: {
		"info": output: dagger.#Logger
		"disabled": {
			output: dagger.#Logger
			level:  dagger.#DisabledLevel
		}
		"warn": {
			output: dagger.#Logger
			level:  dagger.#WarnLevel
		}
	}
	actions: {
		test: {
			helloWorld: core.#Log & {
				input: client.logger."info".output
				message: "hello world 2"
			}

			disabledHelloWorld: core.#Log & {
				input: client.logger."disabled".output
				message: "hello world 3"
			}

			filteredHelloWorld: core.#Log & {
				input: client.logger."warn".output
				message: "hello world 4"
				level: core.#InfoLevel
			}

			moreFields: core.#Log & {
				input: client.logger."info".output
				message: "moreFields message"
				fields: {
					"hello": "world"
					"fun": 22
					"foo": true
				}
			}

			warnLevel: core.#Log & {
				input: client.logger."info".output
				message: "a message"
				fields: {
					"hello": "world"
				}
				level: core.#WarnLevel
			}

			badFields: core.#Log & {
				input: client.logger."info".output
				message: "badFields"
				fields: {
					"badfield": {"a" : "b"}
				}
			}
		}
	}
}
