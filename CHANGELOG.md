# Changelog

## [1.7.1](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.7.0...v1.7.1) (2024-05-29)


### Bug Fixes

* **deps:** bump github.com/go-jose/go-jose/v3 from 3.0.0 to 3.0.3 ([601fcc5](https://github.com/soerenschneider/vault-ssh-cli/commit/601fcc59958441e54946cd3b8474c82d0c7a5f9d))
* **deps:** bump github.com/go-playground/validator/v10 ([83ca68f](https://github.com/soerenschneider/vault-ssh-cli/commit/83ca68fbe8a2261ca71a4443dd34d2f6dd58ae1d))
* **deps:** bump github.com/prometheus/common from 0.44.0 to 0.53.0 ([9f2f585](https://github.com/soerenschneider/vault-ssh-cli/commit/9f2f585e44d5b93e31630b231b19020d079e6b67))
* **deps:** bump github.com/spf13/cobra from 1.7.0 to 1.8.0 ([f64952b](https://github.com/soerenschneider/vault-ssh-cli/commit/f64952be5cad0cd2edf6cec5ca1fa95da0d75e10))
* **deps:** bump golang.org/x/crypto from 0.21.0 to 0.23.0 ([bc838da](https://github.com/soerenschneider/vault-ssh-cli/commit/bc838dac168c890215877a61e5dd8b3d3e189ec0))
* **deps:** bump golang.org/x/net from 0.19.0 to 0.23.0 ([bf3d411](https://github.com/soerenschneider/vault-ssh-cli/commit/bf3d411c19d02826c10e4795800b7278559c9766))
* **deps:** bump golang.org/x/sys from 0.18.0 to 0.20.0 ([c0a5df1](https://github.com/soerenschneider/vault-ssh-cli/commit/c0a5df19fdf26261fd64415ca2b3c82aad076eae))
* **deps:** bump google.golang.org/protobuf from 1.31.0 to 1.33.0 ([7b285e1](https://github.com/soerenschneider/vault-ssh-cli/commit/7b285e1737a15f7b1e01c8c2ecb4465ebdccb469))
* fix dumping metrics with new prometheus version ([47f9b6c](https://github.com/soerenschneider/vault-ssh-cli/commit/47f9b6c8698c9d6ce2c5530cbea56c8dc9f06a12))
* use correct method ([6cdc1a7](https://github.com/soerenschneider/vault-ssh-cli/commit/6cdc1a7613d78fcd6721cdc7886da953fd4629d9))

## [1.7.0](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.6.2...v1.7.0) (2024-05-29)


### Features

* add support for setting ttl and principals ([2e7b6d1](https://github.com/soerenschneider/vault-ssh-cli/commit/2e7b6d1d61560cda540ecc97b6c30e205e0f16b7))


### Bug Fixes

* 121: Provide SshMountPath instead of VaultRole for readCACert: ([b58b26c](https://github.com/soerenschneider/vault-ssh-cli/commit/b58b26c38134853a9ab69cf2065b6c3c53b674b0))
* expand all files ([3781095](https://github.com/soerenschneider/vault-ssh-cli/commit/3781095158d1d39144722a3dd8a7531349f0b829))

## [1.6.2](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.6.1...v1.6.2) (2023-12-24)


### Bug Fixes

* fix print version cmd ([6ff4776](https://github.com/soerenschneider/vault-ssh-cli/commit/6ff4776552e437c0ca7341a8aa81d8257d4c7c0e))

## [1.6.1](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.6.0...v1.6.1) (2023-12-23)


### Bug Fixes

* **deps:** bump github.com/prometheus/client_golang ([#88](https://github.com/soerenschneider/vault-ssh-cli/issues/88)) ([ca303d8](https://github.com/soerenschneider/vault-ssh-cli/commit/ca303d8ffe7692dc975677892aec64c86a9bc64d))
* **deps:** bump github.com/spf13/viper from 1.16.0 to 1.17.0 ([#91](https://github.com/soerenschneider/vault-ssh-cli/issues/91)) ([8c5a196](https://github.com/soerenschneider/vault-ssh-cli/commit/8c5a196d5d88b894119c23dfbb134f02a905992e))
* **deps:** bump golang.org/x/crypto from 0.13.0 to 0.14.0 ([#90](https://github.com/soerenschneider/vault-ssh-cli/issues/90)) ([f826ed4](https://github.com/soerenschneider/vault-ssh-cli/commit/f826ed48f2cd88ad6f903dff747058ddea25c04d))
* **deps:** bump golang.org/x/net from 0.15.0 to 0.17.0 ([#93](https://github.com/soerenschneider/vault-ssh-cli/issues/93)) ([44807d7](https://github.com/soerenschneider/vault-ssh-cli/commit/44807d74dfd4feb61c50b64d3e71a182d5562e2a))
* only dump metrics defined by the tool itself ([a32ecfe](https://github.com/soerenschneider/vault-ssh-cli/commit/a32ecfe23cc58454e8fe2a71f879c9a05076543d))

## [1.6.0](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.5.1...v1.6.0) (2023-09-22)


### Features

* set global loglevel, auto detect whether running on terminal ([beba023](https://github.com/soerenschneider/vault-ssh-cli/commit/beba023a057858abd0be7459668cbb79ddcbc95d))

## [1.5.1](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.5.0...v1.5.1) (2023-08-07)


### Bug Fixes

* expand user paths only at single point ([fe2b311](https://github.com/soerenschneider/vault-ssh-cli/commit/fe2b3114e9b1d4db12ca700b463ded1f16d9022a))

## [1.5.0](https://github.com/soerenschneider/ssh-key-signer/compare/v1.4.1...v1.5.0) (2023-08-02)


### Miscellaneous Chores

* release 1.5.0 ([64ada5a](https://github.com/soerenschneider/ssh-key-signer/commit/64ada5a49845717b0aeb9cd0ddbd1dee30e56207))

## [1.4.1](https://github.com/soerenschneider/ssh-key-signer/compare/v1.4.0...v1.4.1) (2023-05-02)


### Bug Fixes

* respect '--force-new-signature' flag ([8e3a62b](https://github.com/soerenschneider/ssh-key-signer/commit/8e3a62bcd119cff80738fea348b37bbe15f8f0cb))

## [1.4.0](https://github.com/soerenschneider/ssh-key-signer/compare/v1.3.1...v1.4.0) (2023-02-09)


### Features

* sign user certs ([#38](https://github.com/soerenschneider/ssh-key-signer/issues/38)) ([2f06648](https://github.com/soerenschneider/ssh-key-signer/commit/2f066481eb34005a53cc5d5025529f4ad7149572))

## [1.3.1](https://github.com/soerenschneider/ssh-key-signer/compare/v1.3.0...v1.3.1) (2023-01-24)


### Miscellaneous Chores

* release 1.3.1 ([ee4fe98](https://github.com/soerenschneider/ssh-key-signer/commit/ee4fe98c7eca86543c68097a96e03db78861211b))
* release 1.3.3 ([6c234cb](https://github.com/soerenschneider/ssh-key-signer/commit/6c234cb41f0a0ff98a1b44033a9785612473b6d3))

## [1.3.0](https://github.com/soerenschneider/ssh-key-signer/compare/v1.2.1...v1.3.0) (2023-01-24)


### Features

* Add subcmd for reading ca ([#33](https://github.com/soerenschneider/ssh-key-signer/issues/33)) ([8a1caa2](https://github.com/soerenschneider/ssh-key-signer/commit/8a1caa2fa1def937b21ce7fba982733ab3d6218b))

## [1.2.1](https://github.com/soerenschneider/ssh-key-signer/compare/v1.2.0...v1.2.1) (2023-01-06)


### Bug Fixes

* do not write 'go_' prefixed metrics to file ([d270167](https://github.com/soerenschneider/ssh-key-signer/commit/d27016788bbb1028f2f86aaa5273179846efd5bc))

## [1.2.0](https://www.github.com/soerenschneider/ssh-key-signer/compare/v1.1.0...v1.2.0) (2022-02-08)


### Features

* add version info, log on startup ([ce23d40](https://www.github.com/soerenschneider/ssh-key-signer/commit/ce23d40fb5ed7de5a2637718ff16030f69aab4c7))
* split functionality into commands ([1dbdbde](https://www.github.com/soerenschneider/ssh-key-signer/commit/1dbdbde75b1fa88b44ad9f61d0dc93c7e98433e5))

## [1.1.0](https://www.github.com/soerenschneider/ssh-key-signer/compare/v1.0.0...v1.1.0) (2022-02-03)


### Features

* expose new metrics ([79eb199](https://www.github.com/soerenschneider/ssh-key-signer/commit/79eb19948bd3688a61e2f0f0a01e7993c98954bc))

## 1.0.0 (2022-02-02)


### Bug Fixes

* Fix incorrect check for AppRole auth data ([0320284](https://www.github.com/soerenschneider/ssh-key-signer/commit/032028422be048f0c874e23c969db802ea929dd2))
* Include 'cert_type' parameter ([5751678](https://www.github.com/soerenschneider/ssh-key-signer/commit/5751678afef2a69b62048ab4c76ee535748229ec))
* Set AppRole mount path ([96b587f](https://www.github.com/soerenschneider/ssh-key-signer/commit/96b587f3607ee972fb9d150f3fa8cff3bd3ff937))
